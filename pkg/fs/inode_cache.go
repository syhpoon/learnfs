// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"

	"learnfs/pkg/device"
)

type InodeCache struct {
	cache map[InodePtr]*Inode
	dev   device.Device
	pool  *BufPool
	sb    *Superblock
	alloc InodeAllocator

	sync.RWMutex
}

func NewInodeCache(dev device.Device, alloc InodeAllocator,
	pool *BufPool, sb *Superblock) *InodeCache {

	return &InodeCache{
		cache: map[InodePtr]*Inode{},
		pool:  pool,
		sb:    sb,
		dev:   dev,
		alloc: alloc,
	}
}

func (ic *InodeCache) AddInode(inode *Inode) {
	ic.Lock()
	ic.cache[inode.Ptr()] = inode
	ic.Unlock()
}

func (ic *InodeCache) GetCachedInode(ptr InodePtr) *Inode {
	ic.RLock()
	ino := ic.cache[ptr]
	ic.RUnlock()

	return ino
}

func (ic *InodeCache) GetInode(ptr InodePtr) (*Inode, error) {
	ic.RLock()
	ino, ok := ic.cache[ptr]
	ic.RUnlock()

	if ok {
		return ino, nil
	}

	if !ic.alloc.IsAllocated(ptr) {
		return nil, ErrorNotFound
	}

	buf := ic.pool.GetRead()

	op := &device.Op{
		Buf:    buf,
		Offset: ic.sb.InodeOffset(ptr),
	}

	if err := ic.dev.Read(op); err != nil {
		return nil, fmt.Errorf("failed to read inode %d: %w", ptr, err)
	}

	if ino, err := NewInodeFromBuf(ptr, buf); err != nil {
		return nil, err
	} else {
		ic.AddInode(ino)

		return ino, nil
	}
}

func (ic *InodeCache) flush() error {
	ic.Lock()
	defer ic.Unlock()

	var ops []*device.Op

	for ptr, ino := range ic.cache {
		if !ino.IsDirty() {
			continue
		}

		bin, err := serialize(ic.pool, ino)
		if err != nil {
			return fmt.Errorf("failed to serialize inode: %w", err)
		}

		ops = append(ops, &device.Op{
			Buf:    bin,
			Offset: ic.sb.InodeOffset(ptr),
		})

		// TODO: Use the value from device here
		if len(ops) >= 10 {
			if err := ic.dev.Write(ops...); err != nil {
				return fmt.Errorf("failed to write inodes: %w", err)
			}

			ops = ops[:0]
		}
	}

	if len(ops) > 0 {
		if err := ic.dev.Write(ops...); err != nil {
			return fmt.Errorf("failed to write inodes: %w", err)
		}
	}

	return nil
}
