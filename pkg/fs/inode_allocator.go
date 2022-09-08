// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

type InodeAllocator interface {
	AllocateInode() (*Inode, error)
	DeallocateInode(InodePtr) error
	IsAllocated(InodePtr) bool
	GetBitmap() *bitmap
}
