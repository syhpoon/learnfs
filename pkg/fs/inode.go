// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const DIRECT_BLOCKS_NUM = 21

type Inode struct {
	// File type and mode
	Mode uint32

	// Number of directories containing this inode
	Nlink uint32

	// Owner's user id
	Uid uint32

	// Owner's group id
	Gid uint32

	// File size in bytes
	Size uint32

	// Access time
	Atime int64

	// Modification time
	Mtime int64

	// Status change time
	Ctime int64

	// Direct block pointers
	Blocks [DIRECT_BLOCKS_NUM]BlockPtr

	// Single indirect block pointers
	IndirectBlock BlockPtr

	// Double indirect block pointers
	DoubleIndirectBlock BlockPtr

	// Inode number
	ptr InodePtr

	// Dirty flag indicates if inode has been modified
	dirty int32

	sync.RWMutex
}

func newInodeFromBuf(ptr InodePtr, buf Buf) (*Inode, error) {
	ino := &Inode{ptr: ptr}

	err := ino.DecodeFrom(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal inode %d: %w", ptr, err)
	}

	return ino, nil
}

func newInode(ptr InodePtr) *Inode {
	ts := time.Now().UnixMicro()

	return &Inode{
		Mode:                0,
		Nlink:               0,
		Uid:                 0,
		Gid:                 0,
		Size:                0,
		Atime:               ts,
		Mtime:               ts,
		Ctime:               ts,
		Blocks:              [21]BlockPtr{},
		IndirectBlock:       0,
		DoubleIndirectBlock: 0,
		ptr:                 ptr,
		dirty:               0,
	}
}

func (ino *Inode) Ptr() InodePtr {
	return ino.ptr
}

func (ino *Inode) GetBlockPtrs() []BlockPtr {
	var blocks []BlockPtr

	for _, ptr := range ino.Blocks {
		if ptr != 0 {
			blocks = append(blocks, ptr)
		}
	}

	return blocks
}

func (ino *Inode) AddBlockPtr(newPtr BlockPtr) error {
	for i, ptr := range ino.Blocks {
		if ptr == 0 {
			ino.Blocks[i] = newPtr

			return nil
		}
	}

	return fmt.Errorf("all inode blocks are taken")
}

func (ino *Inode) SetDirty(dirty bool) {
	val := int32(1)
	if !dirty {
		val = 0
	}

	atomic.StoreInt32(&ino.dirty, val)
}

func (ino *Inode) IsDirty() bool {
	return atomic.LoadInt32(&ino.dirty) == 1
}

func (ino *Inode) EncodeTo(w io.Writer) error {
	fields := []any{
		ino.Mode,
		ino.Nlink,
		ino.Uid,
		ino.Gid,
		ino.Size,
		ino.Atime,
		ino.Mtime,
		ino.Ctime,
		ino.Blocks,
		ino.IndirectBlock,
		ino.DoubleIndirectBlock,
	}

	return encodeFields(w, fields)
}

func (ino *Inode) DecodeFrom(r io.Reader) error {
	fields := []any{
		&ino.Mode,
		&ino.Nlink,
		&ino.Uid,
		&ino.Gid,
		&ino.Size,
		&ino.Atime,
		&ino.Mtime,
		&ino.Ctime,
		&ino.Blocks,
		&ino.IndirectBlock,
		&ino.DoubleIndirectBlock,
	}

	return decodeFields(r, fields)
}
