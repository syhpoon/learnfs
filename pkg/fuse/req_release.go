// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
)

func (ses *Session) release(ctx context.Context, req *fuse.ReleaseRequest) {
	log.Debug().Interface("req", req).Msg("release")

	handle := uint64(req.Handle)

	ses.Lock()
	if req.Dir {
		_ = ses.dirHandles[handle].stream.CloseSend()
		delete(ses.dirHandles, handle)
	} else {
		delete(ses.fileHandles, handle)
	}
	ses.Unlock()

	req.Respond()
}
