// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) forget(ctx context.Context, req *fuse.ForgetRequest) {
	log.Debug().Interface("req", req).Msg("forget")
	req.Respond()
}
