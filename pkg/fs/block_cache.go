// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"sync"

	"learnfs/pkg/device"
)

type BlockCache struct {
	cache map[BlockPtr]*Block
	dev   device.Device
	pool  *BufPool
	sb    *Superblock
	alloc BlockAllocator

	sync.RWMutex
}

func NewBlockCache(
	dev device.Device, alloc BlockAllocator, pool *BufPool, sb *Superblock) *BlockCache {

	return &BlockCache{
		cache: map[BlockPtr]*Block{},
		pool:  pool,
		sb:    sb,
		dev:   dev,
		alloc: alloc,
	}
}

func (bc *BlockCache) AddBlock(block *Block) {
	bc.RLock()
	bc.cache[block.Ptr()] = block
	bc.RUnlock()
}

func (bc *BlockCache) GetBlock(ptr BlockPtr) (*Block, error) {
	bc.RLock()
	blk, ok := bc.cache[ptr]
	bc.RUnlock()

	if ok {
		return blk, nil
	}

	if !bc.alloc.IsAllocated(ptr) {
		return nil, ErrorNotFound
	}

	buf := bc.pool.GetRead()
	op := &device.Op{
		Buf:    buf,
		Offset: bc.sb.BlockOffset(ptr),
	}

	if err := bc.dev.Read(op); err != nil {
		return nil, fmt.Errorf("failed to read block %d: %w", ptr, err)
	}

	blk = NewBlockFromBuf(ptr, buf)
	bc.AddBlock(blk)

	return blk, nil
}

func (bc *BlockCache) flush() error {
	bc.Lock()
	defer bc.Unlock()

	ops := make([]*device.Op, 0, 10)

	for ptr, blk := range bc.cache {
		if !blk.IsDirty() {
			continue
		}

		ops = append(ops, &device.Op{
			Buf:    blk.Buf(),
			Offset: bc.sb.BlockOffset(ptr),
		})

		// TODO: Use the value from device here
		if len(ops) >= 10 {
			if err := bc.dev.Write(ops...); err != nil {
				return fmt.Errorf("failed to write blocks: %w", err)
			}

			ops = ops[:0]
		}
	}

	if len(ops) > 0 {
		if err := bc.dev.Write(ops...); err != nil {
			return fmt.Errorf("failed to write blocks: %w", err)
		}
	}

	return nil
}
