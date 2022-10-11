// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) readlink(ctx context.Context, req *fuse.ReadlinkRequest) {
	log.Debug().Interface("req", req).Msg("lookup")

	inode := req.Header.Node

	r := &proto.ReadlinkRequest{
		Inode: uint32(inode),
	}

	resp, err := ses.cl.Readlink(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to call Readlink endpoint")
		req.RespondError(err)

		return
	}

	if resp.Errno != 0 {
		req.RespondError(syscall.Errno(resp.Errno))
		return
	}

	req.Respond(resp.Target)
}
