// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"syscall"
)

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
		for _, blkPtr := range inode.getBlockPtrsNoLock() {
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

	buf := make([]byte, 0, size)
	done := false

	for !done && len(buf) < size {
		blockIdx, blockOffset := fs.blockOffset(offset)

		blkPtr, err := fs.getInodeBlock(inode, blockIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to get inode block %d: %w", blockIdx, err)
		}

		// No block allocated for this index yet
		if blkPtr == 0 {
			break
		}

		blk, err := fs.blockCache.getBlock(blkPtr)
		if err != nil {
			return nil, fmt.Errorf("failed to get block %d: %w", ptr, err)
		}

		readSize := size - len(buf)
		sliceSize := int(int64(fs.blockSize) - blockOffset)

		// The required size spans more than one block
		if sliceSize < readSize {
			readSize = sliceSize
		}

		start := blockIdx*int64(fs.blockSize) + blockOffset

		// Avoid exceeding file size
		if start+int64(readSize) > int64(inode.Size) {
			readSize = int(int64(inode.Size) - start)
			done = true
		}

		buf = append(buf, blk.read(int(blockOffset), readSize)...)
	}

	return buf, nil
}

func (fs *Filesystem) Write(ptr InodePtr, offset int64, data []byte) (int, error) {
	inode, err := fs.inodeCache.GetInode(ptr)
	if err != nil {
		return 0, fmt.Errorf("failed to get inode %d: %w", ptr, err)
	}

	inode.Lock()
	defer inode.Unlock()

	origOffset := offset
	var s, e int

	// `data` can be larger than a block size, so we need to split it up
	// and write slices into corresponding blocks
	for e < len(data) {
		var blk *block
		blockIdx, blockOffset := fs.blockOffset(offset)

		// Which actual block back this block index
		blkPtr, err := fs.getInodeBlock(inode, blockIdx)
		if err != nil {
			return 0, fmt.Errorf("failed to get inode block %d: %w", blockIdx, err)
		}

		// Need to allocate a new block
		if blkPtr == 0 {
			blk, err = fs.blockAllocator.AllocateBlock()
			if err != nil {
				return 0, fmt.Errorf("failed to allocate a new block: %w", err)
			}

			if err := fs.setInodeBlock(inode, blockIdx, blk.ptr); err != nil {
				return 0, fmt.Errorf("failed to set inode block %d: %w", blockIdx, err)
			}

			fs.blockCache.addBlock(blk)
		} else {
			blk, err = fs.blockCache.getBlock(ptr)
			if err != nil {
				return 0, fmt.Errorf("failed to get block %d: %w", ptr, err)
			}
		}

		sliceSize := int(int64(fs.blockSize) - blockOffset)
		e += sliceSize
		if e > len(data) {
			e = len(data)
		}

		blk.write(int(blockOffset), data[s:e])
		blk.setDirty(true)
		s += sliceSize
		offset += int64(sliceSize)
	}

	newLen := origOffset + int64(len(data))

	if newLen > int64(inode.Size) {
		inode.Size = uint32(newLen)
	}

	inode.SetDirty(true)

	return len(data), nil
}
