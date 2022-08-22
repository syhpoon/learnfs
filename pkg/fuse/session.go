// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"learnfs/proto"

	"bazil.org/fuse"
	"bazil.org/fuse/fuseutil"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Session struct {
	con        *fuse.Conn
	cl         proto.LearnFSClient
	dir        string
	dirHandle  uint64
	dirHandles map[uint64]*dirHandle

	sync.RWMutex
}

func NewSession(dir string, cl proto.LearnFSClient) (*Session, error) {
	con, err := fuse.Mount(dir,
		fuse.FSName("learnfs"),
		fuse.Subtype("learnfs"),
		fuse.DefaultPermissions(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to mount learnfs: %w", err)
	}

	ses := &Session{
		con:        con,
		dir:        dir,
		cl:         cl,
		dirHandle:  0,
		dirHandles: map[uint64]*dirHandle{},
	}

	return ses, nil
}

func (ses *Session) Start(ctx context.Context) {
	for {
		req, err := ses.con.ReadRequest()
		if err != nil {
			if err == io.EOF {
				return
			}

			log.Error().
				Err(err).
				Msg("failed to read FUSE request")

			continue
		}

		switch r := req.(type) {
		case *fuse.CreateRequest:
			ses.create(ctx, r)
		case *fuse.FlushRequest:
			ses.flush(ctx, r)
		case *fuse.ForgetRequest:
			ses.forget(ctx, r)
		case *fuse.GetattrRequest:
			ses.getAttr(ctx, r)
		case *fuse.GetxattrRequest:
			ses.getxAttr(ctx, r)
		case *fuse.LookupRequest:
			ses.lookup(ctx, r)
		case *fuse.OpenRequest:
			ses.open(ctx, r)
		case *fuse.ReadRequest:
			ses.read(ctx, r)
		case *fuse.ReleaseRequest:
			ses.release(ctx, r)
		case *fuse.SetattrRequest:
			ses.setAttr(ctx, r)
		default:
			fmt.Printf(">>> GOT REQ: %v (%T)\n", req, req)
			req.RespondError(fmt.Errorf("WAT"))
		}
	}
}

func (ses *Session) Stop() {
	_ = fuse.Unmount(ses.dir)
	_ = ses.con.Close()
}

func (ses *Session) getAttr(ctx context.Context, req *fuse.GetattrRequest) {
	log.Debug().Interface("req", req).Msg("getAttr")

	inode := req.Header.Node

	resp, err := ses.cl.GetAttr(ctx, &proto.GetAttrRequest{Inode: uint32(inode)})
	if err != nil {
		log.Error().Err(err).Msg("failed to call GetAttr endpoint")

		req.RespondError(fmt.Errorf("failed to call GetAttr endpoint: %w", err))

		return
	}

	if resp.Errno != 0 {
		req.RespondError(syscall.Errno(resp.Errno))
		return
	}

	req.Respond(&fuse.GetattrResponse{Attr: ses.attr(uint64(inode), resp.Attr)})
}

func (ses *Session) getxAttr(ctx context.Context, req *fuse.GetxattrRequest) {
	log.Debug().Interface("req", req).Msg("getxAttr")
	req.Respond(&fuse.GetxattrResponse{})
}

func (ses *Session) lookup(ctx context.Context, req *fuse.LookupRequest) {
	log.Debug().Interface("req", req).Msg("lookup")

	inode := req.Header.Node

	r := &proto.LookupRequest{
		Inode: uint32(inode),
		Name:  req.Name,
	}

	resp, err := ses.cl.Lookup(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to call Lookup endpoint")
		req.RespondError(err)

		return
	}

	if resp.Errno != 0 {
		req.RespondError(syscall.Errno(resp.Errno))
		return
	}

	req.Respond(&fuse.LookupResponse{
		Node: inode,
		Attr: ses.attr(uint64(inode), resp.Attr),
	})
}

func (ses *Session) open(ctx context.Context, req *fuse.OpenRequest) {
	log.Debug().Interface("req", req).Msg("open")

	inode := req.Header.Node

	stream, err := ses.cl.OpenDir(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to open directory")
		req.RespondError(err)
		return
	}

	if err := stream.Send(&proto.NextDirEntryRequest{Inode: uint32(inode)}); err != nil {
		log.Error().Err(err).Interface("inode", inode).Msg("failed to open directory")
		req.RespondError(err)

		return
	}

	handle := atomic.AddUint64(&ses.dirHandle, 1)

	ses.Lock()
	ses.dirHandles[handle] = &dirHandle{stream: stream}
	ses.Unlock()

	req.Respond(&fuse.OpenResponse{
		Handle: fuse.HandleID(handle),
	})
}

func (ses *Session) read(ctx context.Context, req *fuse.ReadRequest) {
	log.Debug().Interface("req", req).Msg("read")

	handle := uint64(req.Handle)

	ll := log.Error().Uint64("handle", handle)

	if req.Dir {
		h, ok := ses.dirHandles[handle]
		if !ok {
			ll.Msg("invalid directory handle")
			req.RespondError(syscall.EINVAL)
			return
		}

		// Read all the entries
		if h.data == nil {
			stub := &proto.NextDirEntryRequest{}
			for {
				if err := h.stream.Send(stub); err != nil {
					if err == io.EOF {
						break
					}
					ll.Err(err).Msg("failed to request the next dir entry")
					req.RespondError(err)
					return
				}

				resp, err := h.stream.Recv()
				if err != nil {
					ll.Err(err).Msg("failed to receive the next dir entry")
					req.RespondError(err)
					return
				}

				if resp.Errno != 0 {
					req.RespondError(syscall.Errno(resp.Errno))
					return
				}

				if resp.Entry == nil {
					_ = h.stream.CloseSend()
					break
				}

				dirent := fuse.Dirent{
					Inode: uint64(resp.Entry.Inode),
					Type:  ses.entryType(resp.Entry.Mode),
					Name:  resp.Entry.Name,
				}
				h.data = fuse.AppendDirent(h.data, dirent)
			}
		}

		resp := &fuse.ReadResponse{Data: make([]byte, 0, req.Size)}
		fuseutil.HandleRead(req, resp, h.data)
		req.Respond(resp)
	}
}

func (ses *Session) release(ctx context.Context, req *fuse.ReleaseRequest) {
	log.Debug().Interface("req", req).Msg("release")

	handle := uint64(req.Handle)

	if req.Dir {
		ses.Lock()
		_ = ses.dirHandles[handle].stream.CloseSend()
		delete(ses.dirHandles, handle)
		ses.Unlock()

		req.Respond()
	} else {
		// TODO
	}
}

func (ses *Session) create(ctx context.Context, req *fuse.CreateRequest) {
	log.Debug().Interface("req", req).Msg("create")

	dirInode := req.Header.Node

	resp, err := ses.cl.CreateFile(ctx, &proto.CreateFileRequest{
		DirInode: uint32(dirInode),
		Name:     req.Name,
		Flags:    uint32(req.Flags),
		Mode:     ses.fromGoFileMode(req.Mode),
		Umask:    uint32(req.Umask),
		Uid:      req.Uid,
		Gid:      req.Gid,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to create file")
		req.RespondError(err)
		return
	}

	fileInode := fuse.NodeID(resp.Attr.Ino)

	lookup := fuse.LookupResponse{
		Node:       fileInode,
		Attr:       ses.attr(uint64(fileInode), resp.Attr),
		EntryValid: 1 * time.Minute,
	}

	open := fuse.OpenResponse{
		Handle: fuse.HandleID(2),
		Flags:  fuse.OpenKeepCache,
	}

	req.Respond(&fuse.CreateResponse{
		LookupResponse: lookup,
		OpenResponse:   open,
	})
}

func (ses *Session) forget(ctx context.Context, req *fuse.ForgetRequest) {
	log.Debug().Interface("req", req).Msg("forget")
	req.Respond()
}

func (ses *Session) flush(ctx context.Context, req *fuse.FlushRequest) {
	log.Debug().Interface("req", req).Msg("flush")
	req.Respond()
}

func (ses *Session) setAttr(ctx context.Context, req *fuse.SetattrRequest) {
	log.Debug().Interface("req", req).Msg("setAttr")

	inode := req.Header.Node

	r := &proto.SetAttrRequest{
		Inode: uint32(inode),
	}

	if req.Valid.Atime() {
		r.Atime = timestamppb.New(req.Atime)
	}

	if req.Valid.AtimeNow() {
		r.Atime = timestamppb.New(time.Now())
	}

	if req.Valid.Mtime() {
		r.Atime = timestamppb.New(req.Atime)
	}

	if req.Valid.MtimeNow() {
		r.Mtime = timestamppb.New(time.Now())
	}

	if req.Valid.Uid() {
		r.Uid = req.Uid
		r.UidSet = true
	}

	if req.Valid.Gid() {
		r.Gid = req.Gid
		r.GidSet = true
	}

	if req.Valid.Mode() {
		r.Mode = ses.fromGoFileMode(req.Mode)
	}

	resp, err := ses.cl.SetAttr(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to set attributes")
		req.RespondError(err)
		return
	}

	req.Respond(&fuse.SetattrResponse{Attr: ses.attr(uint64(inode), resp.Attr)})
}

func (ses *Session) attr(inode uint64, a *proto.Attr) fuse.Attr {
	return fuse.Attr{
		Valid:     1 * time.Minute,
		Inode:     inode,
		Size:      uint64(a.Size),
		Blocks:    uint64(a.Blocks),
		Atime:     a.Atime.AsTime(),
		Mtime:     a.Mtime.AsTime(),
		Ctime:     a.Ctime.AsTime(),
		Mode:      ses.toGoFileMode(a.Mode),
		Nlink:     a.Nlink,
		Uid:       a.Uid,
		Gid:       a.Gid,
		BlockSize: 4096, // TODO
	}
}

func (ses *Session) entryType(mode uint32) fuse.DirentType {
	switch mode & syscall.S_IFMT {
	case syscall.S_IFSOCK:
		return fuse.DT_Socket
	case syscall.S_IFLNK:
		return fuse.DT_Link
	case syscall.S_IFREG:
		return fuse.DT_File
	case syscall.S_IFBLK:
		return fuse.DT_Block
	case syscall.S_IFDIR:
		return fuse.DT_Dir
	case syscall.S_IFCHR:
		return fuse.DT_Char
	case syscall.S_IFIFO:
		return fuse.DT_FIFO
	default:
		return fuse.DT_Unknown
	}
}

func (ses *Session) toGoFileMode(raw uint32) os.FileMode {
	// Permissions
	mode := os.FileMode(raw & 0777)

	// Type
	switch raw & syscall.S_IFMT {
	case syscall.S_IFBLK:
		mode |= os.ModeDevice
	case syscall.S_IFCHR:
		mode |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		mode |= os.ModeDir
	case syscall.S_IFIFO:
		mode |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		mode |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		mode |= os.ModeSocket
	}

	if raw&syscall.S_ISGID != 0 {
		mode |= os.ModeSetgid
	}
	if raw&syscall.S_ISUID != 0 {
		mode |= os.ModeSetuid
	}
	if raw&syscall.S_ISVTX != 0 {
		mode |= os.ModeSticky
	}

	return mode
}

func (ses *Session) fromGoFileMode(fileMode os.FileMode) uint32 {
	mode := uint32(fileMode.Perm())

	switch fileMode.Type() {
	case 0:
		mode |= syscall.S_IFREG
	case os.ModeDir:
		mode |= syscall.S_IFDIR
	case os.ModeDevice:
		if fileMode&os.ModeCharDevice > 0 {
			mode |= syscall.S_IFCHR
		} else {
			mode |= syscall.S_IFBLK
		}
	case os.ModeNamedPipe:
		mode |= syscall.S_IFIFO
	case os.ModeSymlink:
		mode |= syscall.S_IFLNK
	case os.ModeSocket:
		mode |= syscall.S_IFSOCK
	}

	if fileMode&os.ModeSetuid != 0 {
		mode |= syscall.S_ISUID
	}

	if fileMode&os.ModeSetgid != 0 {
		mode |= syscall.S_ISGID
	}

	return mode
}
