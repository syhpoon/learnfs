// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/syhpoon/learnfs/pkg/device"

	"github.com/rs/zerolog/log"
)

type flusher struct {
	pool           *bufPool
	blockCache     *blockCache
	inodeCache     *InodeCache
	dirCache       *dirCache
	sb             *Superblock
	inodeAllocator InodeAllocator
	blockAllocator BlockAllocator
	dev            device.Device

	sync.Mutex
}

func (f *flusher) flushInode(ptr InodePtr) error {
	f.Lock()
	defer f.Unlock()

	inode := f.inodeCache.GetCachedInode(ptr)
	if inode == nil {
		return nil
	}

	return f.doFlushInode(inode)
}

func (f *flusher) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			f.flushAll()
			f.flushBitmaps()
			return

		case <-time.After(1 * time.Minute):
			f.flushAll()
			f.flushBitmaps()
		}
	}
}

func (f *flusher) flushAll() {
	f.Lock()
	if err := f.inodeCache.flush(); err != nil {
		log.Error().Err(err).Msg("failed to flush inodes")
	}

	if err := f.blockCache.flush(); err != nil {
		log.Error().Err(err).Msg("failed to flush blocks")
	}

	if err := f.dirCache.flush(); err != nil {
		log.Error().Err(err).Msg("failed to flush directories")
	}
	f.Unlock()
}

func (f *flusher) flushBitmaps() {
	f.Lock()
	defer f.Unlock()

	ib := f.inodeAllocator.GetBitmap()
	bb := f.blockAllocator.GetBitmap()

	bin, err := serialize(f.pool, ib.buf, bb.buf)
	if err != nil {
		log.Error().Err(err).Msg("failed to serialize bitmaps")
	}

	op := &device.Op{
		Buf:    bin,
		Offset: uint64(SUPERBLOCK_SIZE),
	}

	if err := f.dev.Write(op); err != nil {
		log.Error().Err(err).Msg("failed to write bitmaps")
	}
}

func (f *flusher) doFlushInode(inode *Inode) error {
	var ops []*device.Op
	var blocks []*block

	// Flush the inode itself
	if inode.IsDirty() {
		bin, err := serialize(f.pool, inode)
		if err != nil {
			return fmt.Errorf("failed to serialize inode: %w", err)
		}

		ops = append(ops, &device.Op{
			Buf:    bin,
			Offset: f.sb.InodeOffset(inode.Ptr()),
		})

		for _, blockPtr := range inode.getBlockPtrs() {
			if block := f.blockCache.getBlockNoFetch(blockPtr); block != nil && block.isDirty() {
				ops = append(ops, &device.Op{
					Buf:    block.data,
					Offset: f.sb.BlockOffset(blockPtr),
				})

				blocks = append(blocks, block)
			}
		}

		if err := f.dev.Write(ops...); err != nil {
			return fmt.Errorf("failed to write ops: %w", err)
		}

		inode.SetDirty(false)
		for i := range blocks {
			blocks[i].setDirty(false)
		}
	}

	return nil
}
