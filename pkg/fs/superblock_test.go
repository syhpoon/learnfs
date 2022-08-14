// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuperblock(t *testing.T) {
	type testCase struct {
		diskSize             uint64
		blockSize            uint32
		inodeSize            uint32
		bytesInodesRatio     uint32
		numInodes            uint32
		numInodeBitmapBlocks uint32
		numDataBitmapBlocks  uint32
		numInodeBlocks       uint32
		firstDataBlockOffset uint64
	}

	cases := []testCase{
		{
			blockSize:        512,
			inodeSize:        64,
			diskSize:         1234567890,
			bytesInodesRatio: 16384,
			// Number of inodes
			// floor(1234567890 / 16384) = floor(75352.044) = 75352
			numInodes: 75352,
			// Number of inode bitmap blocks
			// ceil(75352 / (512 * 8)) = ceil(18.4) = 19
			numInodeBitmapBlocks: 19,
			// Number of data bitmap blocks
			// ceil(2401824 / (512*8)) = ceil(586.38) = 587
			numDataBitmapBlocks: 587,
			// Number of inode blocks
			// ceil(64 * 75352 / 512) = ceil(9419.0) = 9419
			numInodeBlocks:       9419,
			firstDataBlockOffset: 5133824,
		},
	}

	for _, c := range cases {
		params := SuperblockParams{
			DiskSize:         c.diskSize,
			BytesInodesRatio: c.bytesInodesRatio,
			BlockSize:        c.blockSize,
			InodeSize:        c.inodeSize,
		}

		sb := NewSuperblock(params)

		require.Equal(t, SUPERBLOCK_MAGIC, sb.Magic)
		require.Equal(t, c.blockSize, sb.BlockSize)
		require.Equal(t, c.inodeSize, sb.InodeSize)

		// Verify calculated properties
		require.Equal(t, c.numInodes, sb.NumInodes)
		require.Equal(t, c.numInodeBitmapBlocks, sb.NumInodeBitmapBlocks)
		require.Equal(t, c.numDataBitmapBlocks, sb.NumDataBitmapBlocks)
		require.Equal(t, c.numInodeBlocks, sb.NumInodeBlocks)
		require.Equal(t, c.firstDataBlockOffset, sb.FirstDataBlockOffset)
	}
}
