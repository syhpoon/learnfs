// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type block struct {
	ptr   BlockPtr
	data  Buf
	dirty int32

	sync.RWMutex
}

func newBlock(ptr BlockPtr, size uint32) *block {
	return newBlockFromBuf(ptr, make(Buf, size))
}

func newBlockFromBuf(ptr BlockPtr, buf Buf) *block {
	return &block{
		ptr:   ptr,
		data:  buf,
		dirty: 0,
	}
}

func (blk *block) size() int {
	return len(blk.data)
}

func (blk *block) write(offset int, data []byte) {
	blk.Lock()
	copy(blk.data[offset:offset+len(data)], data)
	blk.Unlock()
}

func (blk *block) read(offset int, size int) []byte {
	blk.RLock()
	defer blk.RUnlock()

	buf := make([]byte, size)
	copy(buf, blk.data[offset:size])

	return buf
}

func (blk *block) setDirty(dirty bool) {
	val := int32(1)
	if !dirty {
		val = 0
	}

	atomic.StoreInt32(&blk.dirty, val)
}

func (blk *block) isDirty() bool {
	return atomic.LoadInt32(&blk.dirty) == 1
}

func (blk *block) String() string {
	return fmt.Sprintf("block{ptr=%d, size=%d, dirty=%v}",
		blk.ptr, len(blk.data), blk.dirty)
}
