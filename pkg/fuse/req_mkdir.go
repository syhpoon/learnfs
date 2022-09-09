// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"time"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) mkdir(ctx context.Context, req *fuse.MkdirRequest) {
	log.Debug().Interface("req", req).Msg("mkdir")

	dirInode := req.Header.Node

	resp, err := ses.cl.Mkdir(ctx, &proto.MkdirRequest{
		DirInode: uint32(dirInode),
		Name:     req.Name,
		Mode:     ses.fromGoFileMode(req.Mode),
		Umask:    uint32(req.Umask),
		Uid:      req.Uid,
		Gid:      req.Gid,
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to create directory")
		req.RespondError(err)
		return
	}

	lookup := fuse.LookupResponse{
		Node:       fuse.NodeID(resp.Attr.Ino),
		Attr:       ses.attr(resp.Attr),
		EntryValid: 1 * time.Minute,
	}

	req.Respond(&fuse.MkdirResponse{
		LookupResponse: lookup,
	})
}
