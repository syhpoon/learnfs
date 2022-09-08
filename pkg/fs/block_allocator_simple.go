// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"
)

var _ BlockAllocator = &blockAllocatorSimple{}

type blockAllocatorSimple struct {
	bitmap        *bitmap
	nextFreeBlock *BlockPtr
	blockSize     uint32

	sync.RWMutex
}

func newBlockAllocatorSimple(bm *bitmap, blockSize uint32) *blockAllocatorSimple {
	return &blockAllocatorSimple{
		bitmap:        bm,
		nextFreeBlock: bm.nextClearBit(0),
		blockSize:     blockSize,
	}
}

func (ba *blockAllocatorSimple) GetBitmap() *bitmap {
	return ba.bitmap
}

func (ba *blockAllocatorSimple) AllocateBlock() (*block, error) {
	var ptr BlockPtr

	ba.Lock()
	defer ba.Unlock()

	if ba.nextFreeBlock == nil {
		return nil, fmt.Errorf("no free blocks left")
	} else {
		ptr = *ba.nextFreeBlock
	}

	// Mark the block as used
	ba.bitmap.set(ptr)

	block := newBlock(ptr, ba.blockSize)
	ba.nextFreeBlock = ba.bitmap.nextClearBit(ptr)

	return block, nil
}

func (ba *blockAllocatorSimple) DeallocateBlock(ptr BlockPtr) error {
	ba.Lock()
	defer ba.Unlock()

	ba.bitmap.clear(ptr)

	return nil
}

func (ba *blockAllocatorSimple) IsAllocated(ptr BlockPtr) bool {
	ba.RLock()
	set := ba.bitmap.isSet(ptr)
	ba.RUnlock()

	return set
}
