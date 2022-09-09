// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"fmt"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

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

	req.Respond(&fuse.GetattrResponse{Attr: ses.attr(resp.Attr)})
}
