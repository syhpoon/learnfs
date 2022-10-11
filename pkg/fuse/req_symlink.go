// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"bazil.org/fuse"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/syhpoon/learnfs/proto"
)

func (ses *Session) symlink(ctx context.Context, req *fuse.SymlinkRequest) {
	log.Debug().Interface("req", req).Msg("symlink")

	inode := req.Header.Node

	r := &proto.SymlinkRequest{
		Inode:  uint32(inode),
		Link:   req.NewName,
		Target: req.Target,
		Uid:    req.Uid,
		Gid:    req.Gid,
	}

	resp, err := ses.cl.Symlink(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to symlink")
		req.RespondError(err)
		return
	}

	req.Respond(&fuse.SymlinkResponse{
		LookupResponse: fuse.LookupResponse{
			Node: fuse.NodeID(resp.Attr.Ino),
			Attr: ses.attr(resp.Attr),
		},
	})
}
