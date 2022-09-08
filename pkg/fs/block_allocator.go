// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

type BlockAllocator interface {
	AllocateBlock() (*block, error)
	DeallocateBlock(ptr BlockPtr) error
	IsAllocated(BlockPtr) bool
	GetBitmap() *bitmap
}
