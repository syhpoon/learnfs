// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"
)

var _ InodeAllocator = &inodeAllocatorSimple{}

type inodeAllocatorSimple struct {
	bitmap        *bitmap
	nextFreeInode *InodePtr

	sync.RWMutex
}

func newInodeAllocatorSimple(bm *bitmap) *inodeAllocatorSimple {
	return &inodeAllocatorSimple{
		bitmap:        bm,
		nextFreeInode: bm.nextClearBit(0),
	}
}

func (ia *inodeAllocatorSimple) GetBitmap() *bitmap {
	return ia.bitmap
}

func (ia *inodeAllocatorSimple) AllocateInode() (*Inode, error) {
	var ptr InodePtr

	ia.Lock()
	defer ia.Unlock()

	if ia.nextFreeInode == nil {
		return nil, fmt.Errorf("no free inodes left")
	} else {
		ptr = *ia.nextFreeInode
	}

	// Mark inode as used
	ia.bitmap.set(ptr)

	inode := newInode(ptr)
	ia.nextFreeInode = ia.bitmap.nextClearBit(ptr)

	return inode, nil
}

func (ia *inodeAllocatorSimple) DeallocateInode(ptr InodePtr) error {
	ia.Lock()
	ia.bitmap.clear(ptr)
	ia.Unlock()

	return nil
}

func (ia *inodeAllocatorSimple) IsAllocated(ptr InodePtr) bool {
	ia.RLock()
	set := ia.bitmap.isSet(ptr)
	ia.RUnlock()

	return set
}
