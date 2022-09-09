// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) write(ctx context.Context, req *fuse.WriteRequest) {
	log.Debug().Interface("req", req).Msg("write")

	ses.RLock()
	h, ok := ses.fileHandles[uint64(req.Handle)]
	ses.RUnlock()

	if !ok {
		req.RespondError(syscall.EINVAL)
		return
	}

	resp, err := ses.cl.Write(ctx, &proto.WriteRequest{
		Inode:  uint32(h.inode),
		Offset: req.Offset,
		Data:   req.Data,
	})

	if err != nil {
		log.Error().Uint64("inode", uint64(req.Handle)).Err(err).Msg("failed to write to file")
		req.RespondError(err)
		return
	}

	if resp.Errno != 0 {
		req.RespondError(syscall.Errno(resp.Errno))
		return
	}

	req.Respond(&fuse.WriteResponse{Size: int(resp.Size)})
}
