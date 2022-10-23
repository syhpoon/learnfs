// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"

	"github.com/syhpoon/learnfs/pkg/device"
)

type blockCache struct {
	cache map[BlockPtr]*block
	dev   device.Device
	pool  *bufPool
	sb    *Superblock
	alloc BlockAllocator

	sync.RWMutex
}

func newBlockCache(
	dev device.Device, alloc BlockAllocator, pool *bufPool, sb *Superblock) *blockCache {

	return &blockCache{
		cache: map[BlockPtr]*block{},
		pool:  pool,
		sb:    sb,
		dev:   dev,
		alloc: alloc,
	}
}

func (bc *blockCache) addBlock(block *block) {
	bc.Lock()
	bc.cache[block.ptr] = block
	bc.Unlock()
}

func (bc *blockCache) getBlockNoFetch(ptr BlockPtr) *block {
	bc.RLock()
	block := bc.cache[ptr]
	bc.RUnlock()

	return block
}

func (bc *blockCache) getBlock(ptr BlockPtr) (*block, error) {
	bc.RLock()
	blk, ok := bc.cache[ptr]
	bc.RUnlock()

	if ok {
		return blk, nil
	}

	if !bc.alloc.IsAllocated(ptr) {
		return nil, ErrorNotFound
	}

	buf := bc.pool.getRead()
	op := &device.Op{
		Buf:    buf,
		Offset: bc.sb.BlockOffset(ptr),
	}

	if err := bc.dev.Read(op); err != nil {
		return nil, fmt.Errorf("failed to read block %d: %w", ptr, err)
	}

	blk = newBlockFromBuf(ptr, buf)
	bc.addBlock(blk)

	return blk, nil
}

func (bc *blockCache) flush() error {
	bc.Lock()
	defer bc.Unlock()

	blocks := make([]*block, 0, 10)
	ops := make([]*device.Op, 0, 10)

	for ptr, blk := range bc.cache {
		if !blk.isDirty() {
			continue
		}

		blocks = append(blocks, blk)
		blk.Lock()

		ops = append(ops, &device.Op{
			Buf:    blk.data,
			Offset: bc.sb.BlockOffset(ptr),
		})

		// TODO: Use the value from device here
		if len(ops) >= 10 {
			if err := bc.dev.Write(ops...); err != nil {
				return fmt.Errorf("failed to write blocks: %w", err)
			}

			for _, blk := range blocks {
				blk.setDirty(false)
				blk.Unlock()
			}

			ops = ops[:0]
			blocks = blocks[:0]
		}
	}

	if len(ops) > 0 {
		if err := bc.dev.Write(ops...); err != nil {
			return fmt.Errorf("failed to write blocks: %w", err)
		}

		for _, blk := range blocks {
			blk.setDirty(false)
			blk.Unlock()
		}
	}

	return nil
}
