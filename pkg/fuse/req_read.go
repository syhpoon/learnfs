// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import (
	"context"
	"io"
	"syscall"

	"github.com/syhpoon/learnfs/proto"

	"bazil.org/fuse"
	"bazil.org/fuse/fuseutil"
	"github.com/rs/zerolog/log"
)

func (ses *Session) read(ctx context.Context, req *fuse.ReadRequest) {
	log.Debug().Interface("req", req).Msg("read")

	handle := uint64(req.Handle)

	ll := log.Error().Uint64("handle", handle)

	if req.Dir {
		h, ok := ses.dirHandles[handle]
		if !ok {
			ll.Msg("invalid directory handle")
			req.RespondError(syscall.EINVAL)
			return
		}

		// Read all the entries
		if h.data == nil {
			stub := &proto.NextDirEntryRequest{}
			for {
				if err := h.stream.Send(stub); err != nil {
					if err == io.EOF {
						break
					}
					ll.Err(err).Msg("failed to request the next dir entry")
					req.RespondError(err)
					return
				}

				resp, err := h.stream.Recv()
				if err != nil {
					ll.Err(err).Msg("failed to receive the next dir entry")
					req.RespondError(err)
					return
				}

				if resp.Errno != 0 {
					req.RespondError(syscall.Errno(resp.Errno))
					return
				}

				if resp.Entry == nil {
					_ = h.stream.CloseSend()
					break
				}

				dirent := fuse.Dirent{
					Inode: uint64(resp.Entry.Inode),
					Type:  ses.entryType(resp.Entry.Mode),
					Name:  resp.Entry.Name,
				}
				h.data = fuse.AppendDirent(h.data, dirent)
			}
		}

		resp := &fuse.ReadResponse{Data: make([]byte, 0, req.Size)}
		fuseutil.HandleRead(req, resp, h.data)
		req.Respond(resp)
	} else {
		h, ok := ses.fileHandles[handle]
		if !ok {
			ll.Msg("invalid file handle")
			req.RespondError(syscall.EINVAL)
			return
		}

		resp, err := ses.cl.Read(ctx, &proto.ReadRequest{
			Inode: uint32(h.inode),
			Size:  int32(req.Size),
		})

		if err != nil {
			log.Error().Uint64("inode", uint64(h.inode)).Err(err).Msg("failed to read file")
			req.RespondError(err)
			return
		}

		if resp.Errno != 0 {
			req.RespondError(syscall.Errno(resp.Errno))
			return
		}

		req.Respond(&fuse.ReadResponse{Data: resp.Data})
	}
}
