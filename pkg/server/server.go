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

	"learnfs/pkg/fs"
	"learnfs/proto"

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

var _ proto.LearnFSServer = &Server{}

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

	log.Debug().Interface("req", req).Msg("got GetAttr request")

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

func (s *Server) Lookup(
	ctx context.Context, req *proto.LookupRequest) (*proto.LookupResponse, error) {

	log.Debug().Interface("req", req).Msg("got Lookup request")

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

	log.Debug().Msg("got OpenDir request")

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
				Mode:  uint32(inode.Mode),
			},
		}
		_ = stream.Send(resp)

		idx++
	}
}

func (s *Server) inode2attr(inode *fs.Inode) *proto.Attr {
	attr := &proto.Attr{
		Ino:    inode.Ptr(),
		Size:   inode.Size,
		Blocks: inode.Size * s.fs.BlockSize(),
		Atime:  timestamppb.New(time.UnixMicro(inode.Atime)),
		Mtime:  timestamppb.New(time.UnixMicro(inode.Mtime)),
		Ctime:  timestamppb.New(time.UnixMicro(inode.Ctime)),
		Kind:   nil,
		Mode:   inode.Mode,
		Nlink:  inode.Nlink,
		Uid:    inode.Uid,
		Gid:    inode.Gid,
	}

	switch inode.Mode & syscall.S_IFMT {
	case syscall.S_IFDIR:
		attr.Kind = &proto.Attr_Directory{}
	case syscall.S_IFREG:
		attr.Kind = &proto.Attr_File{}
	case syscall.S_IFLNK:
		attr.Kind = &proto.Attr_Link{}
	case syscall.S_IFSOCK:
		attr.Kind = &proto.Attr_Socket{}
	case syscall.S_IFIFO:
		attr.Kind = &proto.Attr_Fifo{}
	case syscall.S_IFBLK:
		attr.Kind = &proto.Attr_Block{}
	case syscall.S_IFCHR:
		attr.Kind = &proto.Attr_Char{}
	default:
		attr.Kind = &proto.Attr_Unknown{}
	}

	return attr
}
