// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/syhpoon/learnfs/pkg/fs"
	"github.com/syhpoon/learnfs/proto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Params struct {
	Listen string
	Fs     *fs.Filesystem
}

type Server struct {
	fs  *fs.Filesystem
	ctx context.Context
	proto.UnimplementedLearnFSServer
	sync.RWMutex
}

var _ proto.LearnFSServer = (*Server)(nil)

func RunServer(ctx context.Context, params Params) error {
	log.Info().
		Str("listen", params.Listen).
		Msg("starting learnfs server")

	lis, err := net.Listen("tcp", params.Listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	srv := &Server{
		fs:  params.Fs,
		ctx: ctx,
	}

	proto.RegisterLearnFSServer(s, srv)

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	if err := s.Serve(lis); err != nil {
		select {
		case <-ctx.Done():
			return nil
		default:
			return fmt.Errorf("failed to run server: %w", err)
		}
	}

	return nil
}

func (s *Server) GetAttr(
	ctx context.Context, req *proto.GetAttrRequest) (*proto.GetAttrResponse, error) {

	log.Debug().Interface("req", req).Msg("GetAttr")

	inode, err := s.fs.GetInode(req.Inode)

	if err != nil {
		if errors.Is(err, fs.ErrorNotFound) {
			return &proto.GetAttrResponse{Errno: uint32(syscall.ENOENT)}, nil
		}

		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to get inode")

		return nil, status.Errorf(codes.Internal, "failed to get inode")
	}

	return &proto.GetAttrResponse{Attr: s.inode2attr(inode)}, nil
}

func (s *Server) SetAttr(
	ctx context.Context, req *proto.SetAttrRequest) (*proto.SetAttrResponse, error) {

	log.Debug().Interface("req", req).Msg("SetAttr")

	inode, err := s.fs.GetInode(req.Inode)

	if err != nil {
		if errors.Is(err, fs.ErrorNotFound) {
			return &proto.SetAttrResponse{Errno: uint32(syscall.ENOENT)}, nil
		}

		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to get inode")

		return nil, status.Errorf(codes.Internal, "failed to get inode")
	}

	inode.Lock()
	defer inode.Unlock()

	if req.Atime != nil {
		inode.Atime = req.Atime.AsTime().UnixMicro()
	}

	if req.Mtime != nil {
		inode.Mtime = req.Mtime.AsTime().UnixMicro()
	}

	if req.UidSet {
		inode.Uid = req.Uid
	}

	if req.GidSet {
		inode.Gid = req.Gid
	}

	if req.ModeSet {
		inode.Mode = req.Mode
	}

	return &proto.SetAttrResponse{Attr: s.inode2attr(inode)}, nil
}

func (s *Server) Lookup(
	ctx context.Context, req *proto.LookupRequest) (*proto.LookupResponse, error) {

	log.Debug().Interface("req", req).Msg("Lookup")

	inode, err := s.fs.Lookup(req.Inode, req.Name)
	if err != nil {
		if errors.Is(err, fs.ErrorNotFound) {
			return &proto.LookupResponse{Errno: uint32(syscall.ENOENT)}, nil
		}

		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to get inode")
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.LookupResponse{Attr: s.inode2attr(inode)}, nil
}

func (s *Server) OpenDir(stream proto.LearnFS_OpenDirServer) error {
	var dir *fs.Directory
	var entries []*fs.DirEntry
	idx := 0

	log.Debug().Msg("OpenDir")

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to receive client message")
			return err
		}

		if dir == nil {
			dir, err = s.fs.GetDir(in.Inode)
			if err != nil {
				if err == fs.ErrorNotFound {
					resp := &proto.DirEntryResponse{Errno: uint32(syscall.ENOENT)}

					_ = stream.Send(resp)
					return nil
				}

				log.Error().Err(err).Msg("failed to get directory")
				return err
			}

			entries = dir.GetEntries()
		}

		if idx >= len(entries) || entries[idx] == nil {
			resp := &proto.DirEntryResponse{}
			_ = stream.Send(resp)
			return nil
		}

		entry := entries[idx]

		inode, err := s.fs.GetInode(entry.InodePtr)
		if err != nil {
			e := fmt.Errorf("failed to get entry inode %d: %w\n", entry.InodePtr, err)
			log.Error().Err(e)

			return err
		}

		resp := &proto.DirEntryResponse{
			Entry: &proto.DirEntryResponse_DirEntry{
				Inode: entry.InodePtr,
				Name:  entry.GetName(),
				Mode:  inode.Mode,
			},
		}
		_ = stream.Send(resp)

		idx++
	}
}

func (s *Server) OpenFile(
	ctx context.Context, req *proto.OpenFileRequest) (*proto.OpenFileResponse, error) {

	log.Debug().Interface("req", req).Msg("OpenFile")

	_, err := s.fs.GetInode(req.Inode)
	if err != nil {
		if errors.Is(err, fs.ErrorNotFound) {
			return &proto.OpenFileResponse{Errno: uint32(syscall.ENOENT)}, nil
		}

		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to get inode")
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.OpenFileResponse{}, nil
}

func (s *Server) CreateFile(
	ctx context.Context, req *proto.CreateFileRequest) (*proto.CreateFileResponse, error) {

	log.Debug().Interface("req", req).Msg("CreateFile")

	inode, err := s.fs.CreateFile(
		req.DirInode,
		req.Name,
		req.Mode,
		req.Umask,
		req.Uid,
		req.Gid)
	if errors.Is(err, fs.ErrorAlreadyExists) {
		return &proto.CreateFileResponse{Errno: uint32(syscall.EEXIST)}, nil
	}

	if err != nil {
		log.Error().Err(err).
			Uint32("ptr", req.DirInode).
			Str("name", req.Name).
			Msg("failed to create file")

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.CreateFileResponse{Attr: s.inode2attr(inode)}, nil
}

func (s *Server) Flush(ctx context.Context, req *proto.FlushRequest) (*proto.FlushResponse, error) {
	log.Debug().Interface("req", req).Msg("Flush")

	if err := s.fs.Flush(req.Inode); err != nil {
		log.Error().Err(err).
			Uint32("ptr", req.Inode).
			Msg("failed to flush file")

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.FlushResponse{}, nil
}

func (s *Server) Mkdir(
	ctx context.Context, req *proto.MkdirRequest) (*proto.MkdirResponse, error) {

	log.Debug().Interface("req", req).Msg("Mkdir")

	inode, err := s.fs.CreateDirectory(
		req.DirInode,
		req.Name,
		req.Mode,
		req.Umask,
		req.Uid,
		req.Gid)
	if errors.Is(err, fs.ErrorAlreadyExists) {
		return &proto.MkdirResponse{Errno: uint32(syscall.EEXIST)}, nil
	}

	if err != nil {
		log.Error().Err(err).
			Uint32("ptr", req.DirInode).
			Str("name", req.Name).
			Msg("failed to create directory")

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.MkdirResponse{Attr: s.inode2attr(inode)}, nil
}

func (s *Server) Read(
	ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {

	log.Debug().Interface("req", req).Msg("Read")

	_, err := s.fs.GetInode(req.Inode)
	if err != nil {
		if errors.Is(err, fs.ErrorNotFound) {
			return &proto.ReadResponse{Errno: uint32(syscall.ENOENT)}, nil
		}

		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to get inode")
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	data, err := s.fs.Read(req.Inode, req.Offset, int(req.Size))
	if err != nil {
		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to read file data")
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.ReadResponse{Data: data}, nil
}

func (s *Server) RemoveFile(
	ctx context.Context, req *proto.RemoveFileRequest) (*proto.RemoveFileResponse, error) {

	log.Debug().Interface("req", req).Msg("RemoveFile")

	err := s.fs.RemoveFile(req.DirInode, req.Name)
	if errors.Is(err, fs.ErrorNotFound) {
		return &proto.RemoveFileResponse{Errno: uint32(syscall.ENOENT)}, nil
	}

	if err != nil {
		log.Error().Err(err).
			Uint32("ptr", req.DirInode).
			Str("name", req.Name).
			Msg("failed to remove file")

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.RemoveFileResponse{}, nil
}

func (s *Server) Write(
	ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {

	log.Debug().
		Uint32("inode", req.Inode).
		Int64("offset", req.Offset).
		Int("data-len", len(req.Data)).
		Msg("Write")

	_, err := s.fs.GetInode(req.Inode)
	if err != nil {
		if errors.Is(err, fs.ErrorNotFound) {
			return &proto.WriteResponse{Errno: uint32(syscall.ENOENT)}, nil
		}

		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to get inode")
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	size, err := s.fs.Write(req.Inode, req.Offset, req.Data)
	if err != nil {
		log.Error().Err(err).Uint32("ptr", req.Inode).Msg("failed to write file data")
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &proto.WriteResponse{Size: int32(size)}, nil
}

func (s *Server) inode2attr(inode *fs.Inode) *proto.Attr {
	attr := &proto.Attr{
		Ino:    inode.Ptr(),
		Size:   inode.Size,
		Blocks: inode.Size * s.fs.BlockSize(),
		Atime:  timestamppb.New(time.UnixMicro(inode.Atime)),
		Mtime:  timestamppb.New(time.UnixMicro(inode.Mtime)),
		Ctime:  timestamppb.New(time.UnixMicro(inode.Ctime)),
		Mode:   inode.Mode,
		Nlink:  inode.Nlink,
		Uid:    inode.Uid,
		Gid:    inode.Gid,
	}

	return attr
}
