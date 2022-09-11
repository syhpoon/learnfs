// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"time"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

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

	inode := fuse.NodeID(resp.Attr.Ino)

	lookup := fuse.LookupResponse{
		Node:       inode,
		Attr:       ses.attr(resp.Attr),
		EntryValid: 1 * time.Minute,
	}

	open := fuse.OpenResponse{
		Handle: ses.newFileHandle(inode),
		Flags:  fuse.OpenKeepCache,
	}

	req.Respond(&fuse.CreateResponse{
		LookupResponse: lookup,
		OpenResponse:   open,
	})
}
