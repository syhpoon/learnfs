// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

const (
	SUPERBLOCK_MAGIC uint64 = 0x0073666e7261656c
	SUPERBLOCK_SIZE  uint32 = 1024 // Superblock size
)

type SuperblockParams struct {
	// Disk size in bytes
	DiskSize uint64

	// Bytes/inode ratio. One inode will be created for every <value> bytes
	BytesInodesRatio uint32

	// FS block size
	BlockSize uint32

	// Inode size
	InodeSize uint32
}

// Return default set of superblock parameters
func DefaultSuperblockParams() SuperblockParams {
	return SuperblockParams{
		DiskSize:         0,
		BytesInodesRatio: 16384,
		BlockSize:        4096,
		InodeSize:        128,
	}
}

type Superblock struct {
	// Contains SUPERBLOCK_MAGIC value
	Magic uint64

	// FS block size
	BlockSize uint32

	// Total number of inodes
	NumInodes uint32

	// Size of the inode
	InodeSize uint32

	// Number of inode bitmap blocks
	NumInodeBitmapBlocks uint32

	// Number of data bitmap blocks
	NumDataBitmapBlocks uint32

	// Number of inode blocks
	NumInodeBlocks uint32

	// The offset of the first block containing inode in bytes
	FirstInodeBlockOffset uint64

	// The offset of the first block containing data in bytes
	FirstDataBlockOffset uint64
}

// Create a new superblock instance
func NewSuperblock(params SuperblockParams) *Superblock {
	blockBits := float64(params.BlockSize) * 8

	// Calculate the number of inodes
	// For every `params.BytesInodeRatio` bytes of the disk capacity
	// we need to provision one inode
	numInodes := uint32(math.Floor(
		float64(params.DiskSize) / float64(params.BytesInodesRatio)))

	// How many bitmap blocks we need to track all the required inodes
	// One bitmap block can track `blockSize * 8` inodes
	numInodeBitmapBlocks := uint32(math.Ceil(float64(numInodes) / blockBits))

	// How many blocks required for inodes themselves
	numInodeBlocks := uint32(
		math.Ceil(float64(params.InodeSize*numInodes) / float64(params.BlockSize)))

	// The number of blocks occupied by non-data stuff
	nonDataBlocks := uint32(math.Ceil(float64(SUPERBLOCK_SIZE)/float64(params.BlockSize))) +
		numInodeBitmapBlocks + numInodeBlocks

	// This is the total number of blocks disk has
	totalBlocks := uint32(
		math.Floor(float64(params.DiskSize) / float64(params.BlockSize)))

	// How many actual data blocks are left
	numDataBlocks := totalBlocks - nonDataBlocks

	// How many data bitmap blocks are needed to map data blocks
	numDataBitmapBlocks := uint32(math.Ceil(float64(numDataBlocks) / blockBits))

	firstInodeBlockOffset := uint64(
		SUPERBLOCK_SIZE +
			(numInodeBitmapBlocks * params.BlockSize) +
			(numDataBitmapBlocks * params.BlockSize),
	)

	firstDataBlockOffset := firstInodeBlockOffset + uint64(numInodeBlocks*params.BlockSize)

	return &Superblock{
		Magic:                 SUPERBLOCK_MAGIC,
		BlockSize:             params.BlockSize,
		InodeSize:             params.InodeSize,
		NumInodes:             numInodes,
		NumInodeBitmapBlocks:  numInodeBitmapBlocks,
		NumDataBitmapBlocks:   numDataBitmapBlocks,
		NumInodeBlocks:        numInodeBlocks,
		FirstInodeBlockOffset: firstInodeBlockOffset,
		FirstDataBlockOffset:  firstDataBlockOffset,
	}
}

func (sb *Superblock) InodeBitmapSize() uint32 {
	return sb.NumInodeBitmapBlocks * sb.BlockSize
}

func (sb *Superblock) DataBitmapSize() uint32 {
	return sb.NumDataBitmapBlocks * sb.BlockSize
}

func (sb *Superblock) BlockOffset(ptr BlockPtr) uint64 {
	if ptr == 0 {
		panic("block ptr is 0")
	}

	return sb.FirstDataBlockOffset + uint64(ptr*sb.BlockSize)
}

func (sb *Superblock) InodeOffset(ptr InodePtr) uint64 {
	if ptr == 0 {
		panic("inode ptr is 0")
	}

	return sb.FirstInodeBlockOffset + uint64(ptr*sb.BlockSize)
}

func (sb *Superblock) EncodeTo(w io.Writer) error {
	fields := []any{
		sb.Magic,
		sb.BlockSize,
		sb.NumInodes,
		sb.InodeSize,
		sb.NumInodeBitmapBlocks,
		sb.NumDataBitmapBlocks,
		sb.NumInodeBlocks,
		sb.FirstInodeBlockOffset,
		sb.FirstDataBlockOffset,
	}

	return encodeFields(w, fields)
}

func (sb *Superblock) DecodeFrom(r io.Reader) error {
	fields := []any{
		&sb.Magic,
		&sb.BlockSize,
		&sb.NumInodes,
		&sb.InodeSize,
		&sb.NumInodeBitmapBlocks,
		&sb.NumDataBitmapBlocks,
		&sb.NumInodeBlocks,
		&sb.FirstInodeBlockOffset,
		&sb.FirstDataBlockOffset,
	}

	return decodeFields(r, fields)
}

// Load superblock from raw bytes
func LoadSuperblock(buf Buf) (*Superblock, error) {
	sb := &Superblock{}
	if err := sb.DecodeFrom(bytes.NewReader(buf)); err != nil {
		return nil, fmt.Errorf("failed to unmarshal superblock: %v", err)
	}

	if sb.Magic != SUPERBLOCK_MAGIC {
		return nil, fmt.Errorf("invalid superblock magic: %d", sb.Magic)
	}

	return sb, nil
}
