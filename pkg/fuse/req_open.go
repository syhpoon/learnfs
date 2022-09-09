// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"sync/atomic"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) open(ctx context.Context, req *fuse.OpenRequest) {
	log.Debug().Interface("req", req).Msg("open")

	inode := req.Header.Node

	if req.Dir {
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
			Flags:  fuse.OpenKeepCache,
		})
	} else {
		resp, err := ses.cl.OpenFile(ctx, &proto.OpenFileRequest{
			Inode: uint32(inode),
		})

		if err != nil {
			log.Error().Uint64("inode", uint64(inode)).Err(err).Msg("failed to get file")
			req.RespondError(err)
			return
		}

		if resp.Errno != 0 {
			req.RespondError(syscall.Errno(resp.Errno))
			return
		}

		handle := atomic.AddUint64(&ses.fileHandle, 1)

		ses.Lock()
		ses.fileHandles[handle] = &fileHandle{inode: inode}
		ses.Unlock()

		req.Respond(&fuse.OpenResponse{
			Handle: fuse.HandleID(handle),
			Flags:  fuse.OpenKeepCache,
		})
	}
}
