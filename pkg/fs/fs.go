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

func (fs *Filesystem) CreateFile(
	dirInodePtr InodePtr,
	name string,
	mode uint32,
	umask uint32,
	uid uint32,
	gid uint32) (*Inode, error) {
	// Get the dir by the inode ptr
	dir, err := fs.dirCache.getDir(dirInodePtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get dir for inode %d: %w", dirInodePtr, err)
	}

	// Allocate a new inode
	ino, err := fs.inodeAllocator.AllocateInode()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate inode: %w", err)
	}

	ino.Mode = syscall.S_IFREG | mode
	ino.Uid = uid
	ino.Gid = gid
	ino.Nlink = 1

	if err := dir.AddEntry(name, ino.Ptr()); err != nil {
		return nil, fmt.Errorf("failed to add directory entry: %w", err)
	}

	ino.SetDirty(true)
	fs.inodeCache.AddInode(ino)

	return ino, nil
}

func (fs *Filesystem) RemoveFile(dirInodePtr InodePtr, name string) error {
	dir, err := fs.dirCache.getDir(dirInodePtr)
	if err != nil {
		return fmt.Errorf("failed to get dir for inode %d: %w", dirInodePtr, err)
	}

	inodePtr, err := dir.GetEntry(name)
	if err != nil {
		return fmt.Errorf("failed to get entry %s: %w", name, err)
	}

	inode, err := fs.GetInode(inodePtr)
	if err != nil {
		return fmt.Errorf("failed to get inode %d: %w", inodePtr, err)
	}

	inode.Lock()
	defer inode.Unlock()

	inode.Nlink--

	// Need to remove the file
	if inode.Nlink <= 0 {
		// Deallocate all the data blocks
		for _, blkPtr := range inode.GetBlockPtrs() {
			if err := fs.blockAllocator.DeallocateBlock(blkPtr); err != nil {
				return fmt.Errorf("failed to deallocate block %d: %w", blkPtr, err)
			}
		}

		if err := dir.DeleteEntry(name); err != nil {
			return fmt.Errorf("failed to delete directory entry: %w", err)
		}

		if err := fs.inodeAllocator.DeallocateInode(inode.ptr); err != nil {
			return fmt.Errorf("failed to deallocate inode %d: %w", inode.ptr, err)
		}

		fs.inodeCache.DeleteInode(inodePtr)
	}

	return nil
}

func (fs *Filesystem) CreateDirectory(
	dirInodePtr InodePtr,
	name string,
	mode uint32,
	umask uint32,
	uid uint32,
	gid uint32) (*Inode, error) {

	// Get the dir by the inode ptr
	dir, err := fs.dirCache.getDir(dirInodePtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get dir for inode %d: %w", dirInodePtr, err)
	}

	ino, err := fs.inodeAllocator.AllocateInode()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate inode: %w", err)
	}

	ino.Mode = mode
	ino.Uid = uid
	ino.Gid = gid

	blk, err := fs.blockAllocator.AllocateBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate block: %w", err)
	}

	if err := ino.AddBlockPtr(blk.ptr); err != nil {
		return nil, fmt.Errorf("failed to add block to inode: %w", err)
	}

	if err := dir.AddEntry(name, ino.ptr); err != nil {
		return nil, fmt.Errorf("failed to add directory entry: %w", err)
	}

	ino.Nlink = 1
	ino.SetDirty(true)
	fs.inodeCache.AddInode(ino)

	blk.setDirty(true)
	fs.blockCache.addBlock(blk)

	return ino, nil
}

func (fs *Filesystem) GetInode(ptr InodePtr) (*Inode, error) {
	return fs.inodeCache.GetInode(ptr)
}

func (fs *Filesystem) GetDir(ptr InodePtr) (*Directory, error) {
	return fs.dirCache.getDir(ptr)
}

func (fs *Filesystem) Flush(ptr InodePtr) error {
	return fs.flusher.flushInode(ptr)
}

func (fs *Filesystem) Lookup(ptr InodePtr, name string) (*Inode, error) {
	dir, err := fs.dirCache.getDir(ptr)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory %d: %w", ptr, err)
	}

	entryPtr, err := dir.GetEntry(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get dir entry %s: %w", name, err)
	}

	inode, err := fs.inodeCache.GetInode(entryPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get dir entry inode %d: %w", entryPtr, err)
	}

	return inode, nil
}

func (fs *Filesystem) Read(ptr InodePtr, offset int64, size int) ([]byte, error) {
	inode, err := fs.inodeCache.GetInode(ptr)
	if err != nil {
		return nil, fmt.Errorf("failed to get inode %d: %w", ptr, err)
	}

	inode.RLock()
	defer inode.RUnlock()

	blockIdx, blockOffset := fs.blockOffset(offset)

	if inode.Blocks[blockIdx] == 0 {
		return nil, nil
	} else {
		start := blockIdx*int64(fs.blockSize) + blockOffset
		if start+int64(size) > int64(inode.Size) {
			size = int(int64(inode.Size) - offset)
		}

		blk, err := fs.blockCache.getBlock(inode.Blocks[blockIdx])
		if err != nil {
			return nil, fmt.Errorf("failed to get block %d: %w", inode.Blocks[blockIdx], err)
		}

		return blk.read(int(blockOffset), size), nil
	}
}

func (fs *Filesystem) Write(ptr InodePtr, offset int64, data []byte) (int, error) {
	inode, err := fs.inodeCache.GetInode(ptr)
	if err != nil {
		return 0, fmt.Errorf("failed to get inode %d: %w", ptr, err)
	}

	inode.Lock()
	defer inode.Unlock()

	blockIdx, blockOffset := fs.blockOffset(offset)

	var blk *block

	// Need to allocate a new block
	if inode.Blocks[blockIdx] == 0 {
		blk, err = fs.blockAllocator.AllocateBlock()
		if err != nil {
			return 0, fmt.Errorf("failed to allocate a new block: %w", err)
		}

		inode.Blocks[blockIdx] = blk.ptr
	} else {
		blk, err = fs.blockCache.getBlock(inode.Blocks[blockIdx])
		if err != nil {
			return 0, fmt.Errorf("failed to get block %d: %w", inode.Blocks[blockIdx], err)
		}
	}

	// TODO: need to handle the case where data spans the block boundary
	blk.write(int(blockOffset), data)

	size := len(data)
	start := blockIdx*int64(fs.blockSize) + blockOffset
	newLen := start + int64(size)

	if newLen > int64(inode.Size) {
		inode.Size = uint32(newLen)
	}

	blk.setDirty(true)
	inode.SetDirty(true)
	fs.blockCache.addBlock(blk)

	return len(data), nil
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
