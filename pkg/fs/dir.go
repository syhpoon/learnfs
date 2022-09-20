// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/syhpoon/learnfs/pkg/device"
)

type Directory struct {
	inode      *Inode
	name2inode map[string]InodePtr
	entries    map[BlockPtr][]*DirEntry
	cache      *blockCache
	alloc      BlockAllocator
	pool       *bufPool
	blockSize  uint32
	dirty      int32

	sync.RWMutex
}

func NewDirectory(inode *Inode, cache *blockCache, pool *bufPool, blockSize uint32) *Directory {
	entries := make(map[BlockPtr][]*DirEntry)

	// For every inode block create a slice of nil dir entries to represent
	for _, blkPtr := range inode.getBlockPtrs() {
		entries[blkPtr] = make([]*DirEntry, blockSize/DIR_ENTRY_SIZE)
	}

	dir := &Directory{
		inode:      inode,
		name2inode: make(map[string]InodePtr),
		entries:    entries,
		cache:      cache,
		blockSize:  blockSize,
		pool:       pool,
	}

	return dir
}

func LoadDirectory(inode *Inode, cache *blockCache,
	pool *bufPool, blockSize uint32) (*Directory, error) {

	dir := NewDirectory(inode, cache, pool, blockSize)

	// For every block load the corresponding dir entries
	for _, blkPtr := range inode.getBlockPtrs() {
		blk, err := cache.getBlock(blkPtr)
		if err != nil {
			return nil, fmt.Errorf("failed to load block %d: %w", blkPtr, err)
		}

		dir.entries[blkPtr] = make([]*DirEntry, dir.blockSize/DIR_ENTRY_SIZE)

		// Load all the dir entries in the block
		start := 0
		end := DIR_ENTRY_SIZE

		blk.Lock()
		buf := blk.data

		entryIdx := 0
		for end < blk.size() {
			slice := buf[start:end]
			entry := &DirEntry{}

			// TODO: add direntry magic and check it here
			if err := entry.DecodeFrom(bytes.NewReader(slice)); err != nil {
				blk.Unlock()
				return nil, fmt.Errorf("failed to load dir entry: %w", err)
			}

			if entry.InodePtr == 0 {
				dir.entries[blkPtr][entryIdx] = nil
			} else {
				dir.entries[blkPtr][entryIdx] = entry
				dir.name2inode[entry.GetName()] = entry.InodePtr
			}

			start += DIR_ENTRY_SIZE
			end += DIR_ENTRY_SIZE
			entryIdx++
		}

		blk.Unlock()
	}

	return dir, nil
}

func (dir *Directory) AddEntry(name string, inodePtr InodePtr) error {
	dir.Lock()
	defer dir.Unlock()

	if _, ok := dir.name2inode[name]; ok {
		return ErrorAlreadyExists
	}

	entry, err := NewDirEntry(inodePtr, name)
	if err != nil {
		return fmt.Errorf("failed to create new dir entry: %w", err)
	}

	dir.name2inode[name] = inodePtr

	for _, blkPtr := range dir.inode.getBlockPtrs() {
		// Find the first empty slot
		for i := range dir.entries[blkPtr] {
			if dir.entries[blkPtr][i] == nil {
				dir.entries[blkPtr][i] = entry
				dir.SetDirty(true)

				return nil
			}
		}
	}

	return fmt.Errorf("TODO: allocate more directory blocks")
}

func (dir *Directory) flush(dev device.Device, sb *Superblock) error {
	dir.Lock()
	defer dir.Unlock()

	emptyEntry := &DirEntry{}

	write := func(ptr BlockPtr, entries []any) error {
		bin, err := serialize(dir.pool, entries...)
		if err != nil {
			return fmt.Errorf("failed to serialize dir entries: %w", err)
		}

		op := &device.Op{
			Buf:    bin,
			Offset: sb.BlockOffset(ptr),
		}

		if err := dev.Write(op); err != nil {
			return fmt.Errorf("failed to write dir entries: %w", err)
		}

		return nil
	}

	var entries []any
	for _, blkPtr := range dir.inode.getBlockPtrs() {
		for _, e := range dir.entries[blkPtr] {
			entry := e
			if entry == nil {
				entry = emptyEntry
			}

			entries = append(entries, entry)
		}

		if err := write(blkPtr, entries); err != nil {
			return err
		}

		entries = entries[:0]
	}

	return nil
}

func (dir *Directory) GetEntry(name string) (InodePtr, error) {
	dir.RLock()
	ptr, ok := dir.name2inode[name]
	dir.RUnlock()

	if !ok {
		return 0, ErrorNotFound
	}

	return ptr, nil
}

func (dir *Directory) DeleteEntry(name string) error {
	dir.Lock()
	delete(dir.name2inode, name)

	for _, entries := range dir.entries {
		for j, entry := range entries {
			if entry != nil {
				if entry.GetName() == name {
					entries[j] = nil
				}
			}
		}
	}
	dir.SetDirty(true)
	dir.Unlock()

	return nil
}

func (dir *Directory) GetEntries() []*DirEntry {
	dir.RLock()
	defer dir.RUnlock()

	var entries []*DirEntry

	for blkPtr := range dir.entries {
		for i := range dir.entries[blkPtr] {
			if dir.entries[blkPtr][i] == nil {
				continue
			}

			e := *dir.entries[blkPtr][i]
			entries = append(entries, &e)
		}
	}

	return entries
}

func (dir *Directory) SetDirty(dirty bool) {
	val := int32(1)
	if !dirty {
		val = 0
	}

	atomic.StoreInt32(&dir.dirty, val)
}

func (dir *Directory) IsDirty() bool {
	return atomic.LoadInt32(&dir.dirty) == 1
}
