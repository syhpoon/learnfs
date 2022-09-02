// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"
)

var _ BlockAllocator = &BlockAllocatorSimple{}

type BlockAllocatorSimple struct {
	bitmap        *Bitmap
	nextFreeBlock *BlockPtr
	blockSize     uint32

	sync.RWMutex
}

func NewBlockAllocatorSimple(bm *Bitmap, blockSize uint32) *BlockAllocatorSimple {
	return &BlockAllocatorSimple{
		bitmap:        bm,
		nextFreeBlock: bm.NextClearBit(0),
		blockSize:     blockSize,
	}
}

func (ba *BlockAllocatorSimple) GetBitmap() *Bitmap {
	return ba.bitmap
}

func (ba *BlockAllocatorSimple) AllocateBlock() (*Block, error) {
	var ptr BlockPtr

	ba.Lock()
	defer ba.Unlock()

	if ba.nextFreeBlock == nil {
		return nil, fmt.Errorf("no free blocks left")
	} else {
		ptr = *ba.nextFreeBlock
	}

	// Mark the block as used
	ba.bitmap.Set(ptr)

	block := NewBlock(ptr, ba.blockSize)
	ba.nextFreeBlock = ba.bitmap.NextClearBit(ptr)

	return block, nil
}

func (ba *BlockAllocatorSimple) DeallocateBlock(ptr BlockPtr) error {
	ba.Lock()
	defer ba.Unlock()

	ba.bitmap.Clear(ptr)

	return nil
}

func (ba *BlockAllocatorSimple) IsAllocated(ptr BlockPtr) bool {
	ba.RLock()
	set := ba.bitmap.IsSet(ptr)
	ba.RUnlock()

	return set
}
