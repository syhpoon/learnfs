package fs

import (
	"context"
	"fmt"
	"sync"
	"syscall"

	"github.com/syhpoon/learnfs/pkg/device"
)

type Filesystem struct {
	dev            device.Device
	superblock     *Superblock
	inodeAllocator InodeAllocator
	blockAllocator BlockAllocator
	blockCache     *blockCache
	inodeCache     *InodeCache
	dirCache       *dirCache
	pool           *bufPool
	flusher        *flusher
	blockSize      uint32
}

// Create a new filesystem on a given device
func Create(dev device.Device) error {
	// Create superblock instance
	sbParams := DefaultSuperblockParams()
	sbParams.DiskSize = dev.Capacity()

	pool := newBufPool(sbParams.BlockSize)
	sb := NewSuperblock(sbParams)

	// Write metadata
	if err := writeFsMetadata(pool, sb, dev); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	inodeBitmapBuf := make(Buf, sb.InodeBitmapSize())
	inodeBitmap := newBitmap(sb.NumInodes, inodeBitmapBuf)
	inodeBitmap.set(0)

	dataBitmapBuf := make(Buf, sb.DataBitmapSize())
	dataBitmap := newBitmap(sb.NumDataBlocks, dataBitmapBuf)
	dataBitmap.set(0)

	inodeAllocator := newInodeAllocatorSimple(inodeBitmap)
	blockAllocator := newBlockAllocatorSimple(dataBitmap, sb.BlockSize)
	blockCache := newBlockCache(dev, blockAllocator, pool, sb)
	inodeCache := NewInodeCache(dev, inodeAllocator, pool, sb)
	dirCache := newDirCache(pool, inodeCache, blockCache, sb, dev)
	flusherObj := &flusher{
		pool:           pool,
		blockCache:     blockCache,
		inodeCache:     inodeCache,
		dirCache:       dirCache,
		sb:             sb,
		inodeAllocator: inodeAllocator,
		blockAllocator: blockAllocator,
		dev:            dev,
	}

	fs := &Filesystem{
		dev:            dev,
		superblock:     sb,
		inodeAllocator: inodeAllocator,
		blockAllocator: blockAllocator,
		pool:           pool,
		blockCache:     blockCache,
		inodeCache:     inodeCache,
		dirCache:       dirCache,
		blockSize:      sb.BlockSize,
		flusher:        flusherObj,
	}

	// Create a root directory
	if err := fs.createRootDir(); err != nil {
		return fmt.Errorf("failed to create root directory: %w", err)
	}

	return nil
}

// Load filesystem from a device
func Load(dev device.Device) (*Filesystem, error) {
	offset := uint64(0)

	// Load superblock
	buf := make(Buf, SUPERBLOCK_SIZE)
	op := &device.Op{
		Buf:    buf,
		Offset: offset,
	}

	if err := dev.Read(op); err != nil {
		return nil, fmt.Errorf("failed to read superblock: %w", err)
	}

	sb, err := LoadSuperblock(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to load superblock: %w", err)
	}

	offset += uint64(SUPERBLOCK_SIZE)

	// Load bitmaps
	inodeBitmapBuf := make(Buf, sb.InodeBitmapSize())
	ibmOp := &device.Op{
		Buf:    inodeBitmapBuf,
		Offset: offset,
	}

	dataBitmapBuf := make(Buf, sb.DataBitmapSize())
	dbmOp := &device.Op{
		Buf:    dataBitmapBuf,
		Offset: offset + uint64(sb.InodeBitmapSize()),
	}

	if err := dev.Read(ibmOp, dbmOp); err != nil {
		return nil, fmt.Errorf("failed to read bitmaps: %w", err)
	}

	inodeBitmap := newBitmap(sb.NumInodes, inodeBitmapBuf)
	dataBitmap := newBitmap(sb.NumDataBlocks, dataBitmapBuf)
	pool := newBufPool(sb.BlockSize)

	inodeAllocator := newInodeAllocatorSimple(inodeBitmap)
	blockAllocator := newBlockAllocatorSimple(dataBitmap, sb.BlockSize)
	blockCache := newBlockCache(dev, blockAllocator, pool, sb)
	inodeCache := NewInodeCache(dev, inodeAllocator, pool, sb)
	dirCache := newDirCache(pool, inodeCache, blockCache, sb, dev)
	flusherObj := &flusher{
		pool:           pool,
		blockCache:     blockCache,
		inodeCache:     inodeCache,
		dirCache:       dirCache,
		sb:             sb,
		inodeAllocator: inodeAllocator,
		blockAllocator: blockAllocator,
		dev:            dev,
	}

	return &Filesystem{
		dev:            dev,
		superblock:     sb,
		inodeAllocator: inodeAllocator,
		blockAllocator: blockAllocator,
		pool:           pool,
		blockCache:     blockCache,
		inodeCache:     inodeCache,
		dirCache:       dirCache,
		flusher:        flusherObj,
		blockSize:      sb.BlockSize,
	}, nil
}

func (fs *Filesystem) RunFlusher(ctx context.Context, wg *sync.WaitGroup) {
	fs.flusher.run(ctx, wg)
}

func (fs *Filesystem) Shutdown() error {
	return fs.dev.Close()
}

func (fs *Filesystem) BlockSize() uint32 {
	return fs.superblock.BlockSize
}

func (fs *Filesystem) Superblock() *Superblock {
	return fs.superblock
}

func (fs *Filesystem) createRootDir() error {
	ino, err := fs.inodeAllocator.AllocateInode()
	if err != nil {
		return fmt.Errorf("failed to allocate inode: %w", err)
	}

	ino.Mode = syscall.S_IFDIR | 0777

	block, err := fs.blockAllocator.AllocateBlock()
	if err != nil {
		return fmt.Errorf("failed to allocate block: %w", err)
	}

	if err := ino.AddBlockPtr(block.ptr); err != nil {
		return fmt.Errorf("failed to add block to inode: %w", err)
	}

	dir := NewDirectory(ino, fs.blockCache, fs.pool, fs.blockSize)

	if err := dir.AddEntry(".", ino.Ptr()); err != nil {
		return fmt.Errorf("failed to add directory entry: %w", err)
	}

	if err := dir.AddEntry("..", ino.Ptr()); err != nil {
		return fmt.Errorf("failed to add directory entry: %w", err)
	}

	ino.Nlink = 2
	ino.SetDirty(true)
	fs.inodeCache.AddInode(ino)

	block.setDirty(true)
	fs.blockCache.addBlock(block)

	fs.dirCache.addDir(dir)

	fs.flusher.flushAll()
	fs.flusher.flushBitmaps()

	return nil
}

func writeFsMetadata(pool *bufPool, sb *Superblock, dev device.Device) error {
	inodeBitmap := newBitmap(sb.NumInodes, make([]byte, sb.InodeBitmapSize()))
	inodeBitmap.set(0)

	dataBitmap := newBitmap(sb.NumDataBlocks, make([]byte, sb.DataBitmapSize()))
	dataBitmap.set(0)

	bin, err := serialize(pool, sb, inodeBitmap.buf, dataBitmap.buf)
	if err != nil {
		return fmt.Errorf("failed to serialize metadata: %w", err)
	}

	op := &device.Op{
		Buf:    bin,
		Offset: 0,
	}

	if err := dev.Write(op); err != nil {
		return fmt.Errorf("failed to write metadata blocks: %w", err)
	}

	return nil
}

func (fs *Filesystem) blockOffset(offset int64) (int64, int64) {
	// Get the logical block for the offset
	blockIdx := offset / int64(fs.blockSize)
	// Get the offset within the block
	blockOffset := offset % int64(fs.blockSize)

	return blockIdx, blockOffset
}

func (fs *Filesystem) getInodeBlock(ino *Inode, idx int64) (BlockPtr, error) {
	if err := fs.prepareInodeBlock(ino, idx); err != nil {
		return 0, fmt.Errorf("failed to prepare inode block: %w", err)
	}

	return ino.blocks[idx], nil
}

func (fs *Filesystem) setInodeBlock(ino *Inode, idx int64, ptr BlockPtr) error {
	if err := fs.prepareInodeBlock(ino, idx); err != nil {
		return fmt.Errorf("failed to prepare inode block: %w", err)
	}

	if idx >= DIRECT_BLOCKS_NUM {
		ptrsPerPage := int64(fs.blockSize) / BlockPtrSize
		singleMaxBlocks := int64(DIRECT_BLOCKS_NUM) + ptrsPerPage

		var blk *block
		var s int64
		var err error

		if idx < singleMaxBlocks {
			blk, err = fs.blockCache.getBlock(ino.IndirectBlock)
			if err != nil {
				return fmt.Errorf("failed to get indirect block %d for inode %d: %w",
					ino.IndirectBlock, ino.ptr, err)
			}

			s = (idx - DIRECT_BLOCKS_NUM) * BlockPtrSize
		} else {
			doublePtr := ino.doubleBlocks[idx/int64(fs.blockSize)]
			blk, err = fs.blockCache.getBlock(doublePtr)
			if err != nil {
				return fmt.Errorf("failed to get double indirect block %d for inode %d: %w",
					ino.IndirectBlock, ino.ptr, err)
			}

			s = (idx - singleMaxBlocks) * BlockPtrSize
		}

		blk.Lock()
		binEncoding.PutUint32(blk.data[s:s+BlockPtrSize], ptr)
		blk.Unlock()
		blk.setDirty(true)
	}

	ino.blocks[idx] = ptr
	ino.SetDirty(true)

	return nil
}

func (fs *Filesystem) prepareInodeBlock(ino *Inode, idx int64) error {
	// Block already present (either direct or loaded previously)
	if idx < int64(len(ino.blocks)) {
		return nil
	}

	ptrsPerPage := int64(fs.blockSize) / BlockPtrSize

	// Number of single indirect pointers
	singleMaxBlocks := int64(DIRECT_BLOCKS_NUM) + ptrsPerPage

	// Number of double indirect pointers
	doubleMaxBlocks := singleMaxBlocks + ptrsPerPage*ptrsPerPage

	if idx < singleMaxBlocks {
		// Need to allocate the indirect block pointer
		if ino.IndirectBlock == 0 {
			if err := fs.allocateSingleIndirectBlock(ino, ptrsPerPage); err != nil {
				return err
			}
		} else {
			// Load block pointers
			blocks, err := fs.loadIndirectBlocks(ino.IndirectBlock, ptrsPerPage)
			if err != nil {
				return err
			}

			ino.addBlocksNoLock(blocks)
		}

		ino.SetDirty(true)
	} else if idx < doubleMaxBlocks {
		var doubleBlk *block
		var err error

		// Double indirect
		if ino.DoubleIndirectBlock == 0 {
			doubleBlk, err = fs.allocateDoubleIndirectBlock(ino)
			if err != nil {
				return err
			}
		} else {
			doubleBlk, err = fs.blockCache.getBlock(ino.DoubleIndirectBlock)
			if err != nil {
				return fmt.Errorf("failed to load double indirect block %d: %w",
					ino.DoubleIndirectBlock, err)
			}
		}

		doubleBlk.Lock()
		defer doubleBlk.Unlock()

	LOOP:
		for l1 := int64(0); l1 < ptrsPerPage; l1++ {
			var l1Block *block
			l1Offset := l1 * BlockPtrSize
			l1BlockPtr := binEncoding.Uint32(doubleBlk.data[l1Offset : l1Offset+BlockPtrSize])

			if l1BlockPtr == 0 {
				l1Block, err = fs.allocateL1DoubleIndirectBlock(doubleBlk, l1Offset)
				if err != nil {
					return err
				}
			} else {
				l1Block, err = fs.blockCache.getBlock(l1BlockPtr)
				if err != nil {
					return fmt.Errorf("failed to load l1 double indirect block %d: %w",
						l1BlockPtr, err)
				}
			}

			l1Block.Lock()
			defer l1Block.Unlock()

			for l2 := int64(0); l2 < ptrsPerPage; l2++ {
				l2Offset := l2 * BlockPtrSize
				l2BlockPtr := binEncoding.Uint32(l1Block.data[l2Offset : l2Offset+BlockPtrSize])

				if l2BlockPtr == 0 {
					blk, err := fs.blockAllocator.AllocateBlock()
					if err != nil {
						return fmt.Errorf(
							"failed to allocate l2 double indirect block for inode %d: %w",
							ino.ptr, err)
					}

					fs.blockCache.addBlock(blk)
					binEncoding.PutUint32(l1Block.data[l2Offset:l2Offset+BlockPtrSize], blk.ptr)
					l1Block.setDirty(true)

					// Since it's a newly allocated indirect block, all the pointers are
					// unused (0) at this point
					ino.addBlocksNoLock(make([]BlockPtr, ptrsPerPage))
					ino.doubleBlocks[idx/int64(fs.blockSize)] = blk.ptr
				} else {
					blocks, err := fs.loadIndirectBlocks(l2BlockPtr, ptrsPerPage)
					if err != nil {
						return err
					}

					ino.addBlocksNoLock(blocks)
				}

				if int64(len(ino.blocks)) > idx {
					break LOOP
				}
			}
		}
	} else {
		return fmt.Errorf("inode %d block index exceeding maximum file size: %d", ino.ptr, idx)
	}

	return nil
}

func (fs *Filesystem) allocateSingleIndirectBlock(ino *Inode, ptrsPerPage int64) error {
	blk, err := fs.blockAllocator.AllocateBlock()
	if err != nil {
		return fmt.Errorf(
			"failed to allocate indirect block for inode %d: %w", ino.ptr, err)
	}

	fs.blockCache.addBlock(blk)

	// Since it's a newly allocated indirect block, all the pointers are
	// unused (0) at this point
	ino.blocks = append(ino.blocks, make([]BlockPtr, ptrsPerPage)...)
	ino.IndirectBlock = blk.ptr

	return nil
}

func (fs *Filesystem) allocateDoubleIndirectBlock(ino *Inode) (*block, error) {
	blk, err := fs.blockAllocator.AllocateBlock()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to allocate double indirect block for inode %d: %w", ino.ptr, err)
	}

	fs.blockCache.addBlock(blk)
	ino.DoubleIndirectBlock = blk.ptr

	return blk, nil
}

func (fs *Filesystem) allocateL1DoubleIndirectBlock(
	doubleBlk *block, offset int64) (*block, error) {

	blk, err := fs.blockAllocator.AllocateBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate l1 double indirect block: %w", err)
	}

	fs.blockCache.addBlock(blk)
	binEncoding.PutUint32(doubleBlk.data[offset:offset+BlockPtrSize], blk.ptr)
	doubleBlk.setDirty(true)

	return blk, nil
}

func (fs *Filesystem) loadIndirectBlocks(ptr BlockPtr, ptrsPerPage int64) ([]BlockPtr, error) {
	blk, err := fs.blockCache.getBlock(ptr)
	if err != nil {
		return nil, fmt.Errorf("failed to load indirect block %d: %w", ptr, err)
	}

	blocks := make([]BlockPtr, ptrsPerPage)

	for i := range blocks {
		s := int64(i) * BlockPtrSize
		blocks[i] = binEncoding.Uint32(blk.data[s : s+BlockPtrSize])
	}

	return blocks, nil
}
