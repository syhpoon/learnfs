// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"

	"github.com/syhpoon/learnfs/pkg/device"
)

type dirCache struct {
	cache      map[InodePtr]*Directory
	inodeCache *InodeCache
	blockCache *blockCache
	pool       *bufPool
	blockSize  uint32
	dev        device.Device
	sb         *Superblock

	sync.RWMutex
}

func newDirCache(
	pool *bufPool,
	inodeCache *InodeCache,
	blockCache *blockCache,
	sb *Superblock,
	dev device.Device) *dirCache {

	return &dirCache{
		cache:      map[InodePtr]*Directory{},
		pool:       pool,
		inodeCache: inodeCache,
		blockCache: blockCache,
		blockSize:  sb.BlockSize,
		sb:         sb,
		dev:        dev,
	}
}

func (dc *dirCache) addDir(dir *Directory) {
	dc.Lock()
	dc.cache[dir.inode.Ptr()] = dir
	dc.Unlock()
}

func (dc *dirCache) getDir(ptr InodePtr) (*Directory, error) {
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

func (dc *dirCache) flush() error {
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
