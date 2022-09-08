// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Block struct {
	ptr   BlockPtr
	data  Buf
	dirty int32

	sync.RWMutex
}

func newBlock(ptr BlockPtr, size uint32) *Block {
	return newBlockFromBuf(ptr, make(Buf, size))
}

func newBlockFromBuf(ptr BlockPtr, buf Buf) *Block {
	return &Block{
		ptr:   ptr,
		data:  buf,
		dirty: 0,
	}
}

func (blk *Block) size() int {
	return len(blk.data)
}

func (blk *Block) write(offset int, data []byte) {
	blk.Lock()
	copy(blk.data[offset:offset+len(data)], data)
	blk.Unlock()
}

func (blk *Block) read(offset int, size int) []byte {
	blk.RLock()
	defer blk.RUnlock()

	buf := make([]byte, size)
	copy(buf, blk.data[offset:size])

	return buf
}

func (blk *Block) setDirty(dirty bool) {
	val := int32(1)
	if !dirty {
		val = 0
	}

	atomic.StoreInt32(&blk.dirty, val)
}

func (blk *Block) isDirty() bool {
	return atomic.LoadInt32(&blk.dirty) == 1
}

func (blk *Block) String() string {
	return fmt.Sprintf("Block{ptr=%d, size=%d, dirty=%v}",
		blk.ptr, len(blk.data), blk.dirty)
}
