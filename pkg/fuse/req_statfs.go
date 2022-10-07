// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) statfs(ctx context.Context, req *fuse.StatfsRequest) {
	log.Debug().Interface("req", req).Msg("statfs")

	resp, err := ses.cl.Statfs(ctx, &proto.StatfsRequest{})
	if err != nil {
		log.Error().Err(err).Msg("failed to get statfs")
		req.RespondError(err)
		return
	}

	req.Respond(&fuse.StatfsResponse{
		Blocks: resp.TotalBlocks,
		Bfree:  resp.FreeBlocks,
		Bavail: resp.FreeBlocks,
		Files:  resp.TotalInodes,
		Ffree:  resp.FreeInodes,
		Bsize:  resp.BlockSize,
	})
}
