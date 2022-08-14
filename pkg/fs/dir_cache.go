// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"

	"learnfs/pkg/device"
)

type DirCache struct {
	cache      map[InodePtr]*Directory
	inodeCache *InodeCache
	blockCache *BlockCache
	pool       *BufPool
	blockSize  uint32
	dev        device.Device
	sb         *Superblock

	sync.RWMutex
}

func NewDirCache(
	pool *BufPool,
	inodeCache *InodeCache,
	blockCache *BlockCache,
	sb *Superblock,
	dev device.Device) *DirCache {

	return &DirCache{
		cache:      map[InodePtr]*Directory{},
		pool:       pool,
		inodeCache: inodeCache,
		blockCache: blockCache,
		blockSize:  sb.BlockSize,
		sb:         sb,
		dev:        dev,
	}
}

func (dc *DirCache) AddDir(dir *Directory) {
	dc.RLock()
	dc.cache[dir.inode.Ptr()] = dir
	dc.RUnlock()
}

func (dc *DirCache) GetDir(ptr InodePtr) (*Directory, error) {
	dc.RLock()
	d, ok := dc.cache[ptr]
	dc.RUnlock()

	if ok {
		return d, nil
	}

	ino, err := dc.inodeCache.GetInode(ptr)
	if err != nil {
		return nil, fmt.Errorf("failed to get inode %d: %w", ptr, err)
	}

	d, err = LoadDirectory(ino, dc.blockCache, dc.pool, dc.blockSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %d: %w", ptr, err)
	}

	dc.Lock()
	dc.cache[ptr] = d
	dc.Unlock()

	return d, nil
}

func (dc *DirCache) flush() error {
	dc.Lock()
	defer dc.Unlock()

	for ptr, dir := range dc.cache {
		if !dir.IsDirty() {
			continue
		}

		if err := dir.flush(dc.dev, dc.sb); err != nil {
			return fmt.Errorf("failed to flush dir %d: %w", ptr, err)
		}
	}

	return nil
}
