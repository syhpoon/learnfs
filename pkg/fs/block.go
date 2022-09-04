// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"bytes"
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

func NewBlock(ptr BlockPtr, size uint32) *Block {
	return NewBlockFromBuf(ptr, make(Buf, size))
}

func NewBlockFromBuf(ptr BlockPtr, buf Buf) *Block {
	return &Block{
		ptr:   ptr,
		data:  buf,
		dirty: 0,
	}
}

func (blk *Block) Ptr() BlockPtr {
	return blk.ptr
}

func (blk *Block) Size() int {
	return len(blk.data)
}

func (blk *Block) Buf() Buf {
	return blk.data
}

func (blk *Block) write(offset int, data []byte) {
	blk.Lock()
	copy(blk.data[offset:len(data)], data)
	blk.Unlock()
}

func (blk *Block) read(offset int, size int) []byte {
	blk.RLock()
	defer blk.RUnlock()

	idx := bytes.IndexByte(blk.data, 0)
	if idx < 0 {
		return nil
	}

	if idx > size-1 {
		idx = size - 1
	}

	buf := make([]byte, idx+1)
	copy(buf, blk.data[offset:len(buf)])

	return buf
}

func (blk *Block) SetDirty(dirty bool) {
	val := int32(1)
	if !dirty {
		val = 0
	}

	atomic.StoreInt32(&blk.dirty, val)
}

func (blk *Block) IsDirty() bool {
	return atomic.LoadInt32(&blk.dirty) == 1
}

func (blk *Block) String() string {
	return fmt.Sprintf("Block{ptr=%d, size=%d, dirty=%v}",
		blk.ptr, len(blk.data), blk.dirty)
}
