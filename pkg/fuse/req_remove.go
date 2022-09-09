// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) remove(ctx context.Context, req *fuse.RemoveRequest) {
	log.Debug().Interface("req", req).Msg("remove")

	inode := uint32(req.Node)

	if req.Dir {
		// TODO
	} else {
		resp, err := ses.cl.RemoveFile(ctx, &proto.RemoveFileRequest{
			DirInode: inode,
			Name:     req.Name,
		})

		if err != nil {
			log.Error().Str("file", req.Name).Err(err).Msg("failed to remove file")
			req.RespondError(err)
			return
		}

		if resp.Errno != 0 {
			req.RespondError(syscall.Errno(resp.Errno))
			return
		}
	}

	req.Respond()
}
