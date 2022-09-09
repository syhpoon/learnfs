// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"time"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (ses *Session) setAttr(ctx context.Context, req *fuse.SetattrRequest) {
	log.Debug().Interface("req", req).Msg("setAttr")

	inode := req.Header.Node

	r := &proto.SetAttrRequest{
		Inode: uint32(inode),
	}

	if req.Valid.Atime() {
		r.Atime = timestamppb.New(req.Atime)
	}

	if req.Valid.AtimeNow() {
		r.Atime = timestamppb.New(time.Now())
	}

	if req.Valid.Mtime() {
		r.Mtime = timestamppb.New(req.Mtime)
	}

	if req.Valid.MtimeNow() {
		r.Mtime = timestamppb.New(time.Now())
	}

	if req.Valid.Uid() {
		r.Uid = req.Uid
		r.UidSet = true
	}

	if req.Valid.Gid() {
		r.Gid = req.Gid
		r.GidSet = true
	}

	if req.Valid.Mode() {
		r.Mode = ses.fromGoFileMode(req.Mode)
		r.ModeSet = true
	}

	resp, err := ses.cl.SetAttr(ctx, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to set attributes")
		req.RespondError(err)
		return
	}

	req.Respond(&fuse.SetattrResponse{Attr: ses.attr(resp.Attr)})
}
