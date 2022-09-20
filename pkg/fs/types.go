// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"unsafe"
)

// block pointer
type BlockPtr = uint32

// Inode pointer
type InodePtr = uint32

// Common buffer
type Buf = []byte

const (
	BlockPtrSize = int64(unsafe.Sizeof(BlockPtr(0)))
)
