// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

type dirHandle struct {
	stream proto.LearnFS_OpenDirClient
	data   []byte
}

type fileHandle struct {
	inode fuse.NodeID
}

type Session struct {
	con         *fuse.Conn
	cl          proto.LearnFSClient
	dir         string
	dirHandle   uint64
	fileHandle  uint64
	dirHandles  map[uint64]*dirHandle
	fileHandles map[uint64]*fileHandle

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
		con:         con,
		dir:         dir,
		cl:          cl,
		dirHandle:   0,
		fileHandle:  0,
		dirHandles:  map[uint64]*dirHandle{},
		fileHandles: map[uint64]*fileHandle{},
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
		case *fuse.MkdirRequest:
			ses.mkdir(ctx, r)
		case *fuse.OpenRequest:
			ses.open(ctx, r)
		case *fuse.ReadRequest:
			ses.read(ctx, r)
		case *fuse.ReleaseRequest:
			ses.release(ctx, r)
		case *fuse.RemoveRequest:
			ses.remove(ctx, r)
		case *fuse.SetattrRequest:
			ses.setAttr(ctx, r)
		case *fuse.WriteRequest:
			ses.write(ctx, r)
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

func (ses *Session) attr(a *proto.Attr) fuse.Attr {
	return fuse.Attr{
		Valid:     1 * time.Minute,
		Inode:     uint64(a.Ino),
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
