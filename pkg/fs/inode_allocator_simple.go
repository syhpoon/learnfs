// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"
)

var _ InodeAllocator = &InodeAllocatorSimple{}

type InodeAllocatorSimple struct {
	bitmap        *Bitmap
	nextFreeInode *InodePtr

	sync.RWMutex
}

func NewInodeAllocatorSimple(bm *Bitmap) *InodeAllocatorSimple {
	return &InodeAllocatorSimple{
		bitmap:        bm,
		nextFreeInode: bm.NextClearBit(0),
	}
}

func (ia *InodeAllocatorSimple) GetBitmap() *Bitmap {
	return ia.bitmap
}

func (ia *InodeAllocatorSimple) AllocateInode() (*Inode, error) {
	var ptr InodePtr

	ia.Lock()
	defer ia.Unlock()

	if ia.nextFreeInode == nil {
		return nil, fmt.Errorf("no free inodes left")
	} else {
		ptr = *ia.nextFreeInode
	}

	// Mark inode as used
	ia.bitmap.Set(ptr)

	inode := NewInode(ptr)
	ia.nextFreeInode = ia.bitmap.NextClearBit(ptr)

	return inode, nil
}

func (ia *InodeAllocatorSimple) IsAllocated(ptr InodePtr) bool {
	ia.RLock()
	set := ia.bitmap.IsSet(ptr)
	ia.RUnlock()

	return set
}
