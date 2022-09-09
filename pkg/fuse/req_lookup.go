// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

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
		Node: fuse.NodeID(resp.Attr.Ino),
		Attr: ses.attr(resp.Attr),
	})
}
