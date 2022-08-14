// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fuse

import "learnfs/proto"

type dirHandle struct {
	stream proto.LearnFS_OpenDirClient
	data   []byte
}
