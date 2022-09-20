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

	// Single indirect block pointers
	IndirectBlock BlockPtr

	// Double indirect block pointers
	DoubleIndirectBlock BlockPtr

	// Inode number
	ptr InodePtr

	// Dirty flag indicates if inode has been modified
	dirty int32

	// a linear list of all the inode blocks, direct and indirect
	blocks []BlockPtr

	sync.RWMutex
}

func newInodeFromBuf(ptr InodePtr, buf Buf) (*Inode, error) {
	ino := &Inode{
		ptr:    ptr,
		blocks: make([]BlockPtr, DIRECT_BLOCKS_NUM),
	}

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
		IndirectBlock:       0,
		DoubleIndirectBlock: 0,
		ptr:                 ptr,
		dirty:               0,
		blocks:              make([]BlockPtr, DIRECT_BLOCKS_NUM),
	}
}

func (ino *Inode) Ptr() InodePtr {
	return ino.ptr
}

func (ino *Inode) getBlockPtrs() []BlockPtr {
	ino.RLock()
	blocks := ino.getBlockPtrsNoLock()
	ino.RUnlock()

	return blocks
}

func (ino *Inode) getBlockPtrsNoLock() []BlockPtr {
	var blocks []BlockPtr

	for _, ptr := range ino.blocks {
		if ptr != 0 {
			blocks = append(blocks, ptr)
		}
	}

	return blocks
}

func (ino *Inode) AddBlockPtr(newPtr BlockPtr) error {
	ino.Lock()
	defer ino.Unlock()

	for i, ptr := range ino.blocks {
		if ptr == 0 {
			ino.blocks[i] = newPtr

			return nil
		}
	}

	return fmt.Errorf("TODO: all inode blocks are taken")
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
		ino.blocks[:DIRECT_BLOCKS_NUM],
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
		&ino.blocks,
		&ino.IndirectBlock,
		&ino.DoubleIndirectBlock,
	}

	return decodeFields(r, fields)
}
