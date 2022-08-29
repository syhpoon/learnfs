package fs

import (
	"context"
	"fmt"
	"sync"
	"syscall"

	"learnfs/pkg/device"
)

type Filesystem struct {
	dev            device.Device
	superblock     *Superblock
	inodeAllocator InodeAllocator
	blockAllocator BlockAllocator
	blockCache     *BlockCache
	inodeCache     *InodeCache
	dirCache       *DirCache
	pool           *BufPool
	flusher        *flusher
	blockSize      uint32
}

// Create a new filesystem on a given device
func Create(dev device.Device) error {
	// Create superblock instance
	sbParams := DefaultSuperblockParams()
	sbParams.DiskSize = dev.Capacity()

	pool := NewBufPool(sbParams.BlockSize)
	sb := NewSuperblock(sbParams)

	// Write metadata
	if err := writeFsMetadata(pool, sb, dev); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	inodeBitmapBuf := make(Buf, sb.InodeBitmapSize())
	inodeBitmap := NewBitmap(sb.NumInodes, inodeBitmapBuf)
	inodeBitmap.Set(0)

	dataBitmapBuf := make(Buf, sb.DataBitmapSize())
	dataBitmap := NewBitmap(sb.NumDataBlocks, dataBitmapBuf)
	dataBitmap.Set(0)

	inodeAllocator := NewInodeAllocatorSimple(inodeBitmap)
	blockAllocator := NewBlockAllocatorSimple(dataBitmap, sb.BlockSize)
	blockCache := NewBlockCache(dev, blockAllocator, pool, sb)
	inodeCache := NewInodeCache(dev, inodeAllocator, pool, sb)
	dirCache := NewDirCache(pool, inodeCache, blockCache, sb, dev)
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

	inodeBitmap := NewBitmap(sb.NumInodes, inodeBitmapBuf)
	dataBitmap := NewBitmap(sb.NumDataBlocks, dataBitmapBuf)
	pool := NewBufPool(sb.BlockSize)

	inodeAllocator := NewInodeAllocatorSimple(inodeBitmap)
	blockAllocator := NewBlockAllocatorSimple(dataBitmap, sb.BlockSize)
	blockCache := NewBlockCache(dev, blockAllocator, pool, sb)
	inodeCache := NewInodeCache(dev, inodeAllocator, pool, sb)
	dirCache := NewDirCache(pool, inodeCache, blockCache, sb, dev)
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
	dir, err := fs.dirCache.GetDir(dirInodePtr)
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

func (fs *Filesystem) GetInode(ptr InodePtr) (*Inode, error) {
	return fs.inodeCache.GetInode(ptr)
}

func (fs *Filesystem) GetDir(ptr InodePtr) (*Directory, error) {
	return fs.dirCache.GetDir(ptr)
}

func (fs *Filesystem) Flush(ptr InodePtr) error {
	return fs.flusher.flushInode(ptr)
}

func (fs *Filesystem) Lookup(ptr InodePtr, name string) (*Inode, error) {
	dir, err := fs.dirCache.GetDir(ptr)
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

	block.SetDirty(true)
	fs.blockCache.AddBlock(block)

	fs.dirCache.AddDir(dir)

	fs.flusher.flushAll()
	fs.flusher.flushBitmaps()

	return nil
}

func writeFsMetadata(pool *BufPool, sb *Superblock, dev device.Device) error {
	inodeBitmap := NewBitmap(sb.NumInodes, make([]byte, sb.InodeBitmapSize()))
	inodeBitmap.Set(0)

	dataBitmap := NewBitmap(sb.NumDataBlocks, make([]byte, sb.DataBitmapSize()))
	dataBitmap.Set(0)

	bin, err := serialize(pool, sb, inodeBitmap.GetBuf(), dataBitmap.GetBuf())
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
