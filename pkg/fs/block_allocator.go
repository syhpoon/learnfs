// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

type BlockAllocator interface {
	AllocateBlock() (*Block, error)
	IsAllocated(BlockPtr) bool
	GetBitmap() *Bitmap
}
