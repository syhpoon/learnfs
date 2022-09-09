// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) flush(ctx context.Context, req *fuse.FlushRequest) {
	log.Debug().Interface("req", req).Msg("flush")

	_, err := ses.cl.Flush(ctx, &proto.FlushRequest{Inode: uint32(req.Handle)})
	if err != nil {
		log.Error().Err(err).Uint64("inode", uint64(req.Handle)).Msg("failed to flush file")
		req.RespondError(err)
		return
	}

	req.Respond()
}
