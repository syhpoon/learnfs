// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) getxAttr(ctx context.Context, req *fuse.GetxattrRequest) {
	log.Debug().Interface("req", req).Msg("getxAttr")
	req.Respond(&fuse.GetxattrResponse{})
}
