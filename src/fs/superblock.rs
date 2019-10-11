/*
 MIT License

 Copyright (c) 2019 Max Kuznetsov <syhpoon@syhpoon.ca>

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 SOFTWARE.
*/

use serde::{Serialize, Deserialize};
use failure::{Error, format_err};

use super::{BOOT_BLOCK_SIZE, SUPERBLOCK_SIZE};
use super::Bitmap;

pub struct Params {
    // Disk size in bytes
    pub disk_size: u64,

    // Bytes/inode ration. One inode will be created for every <value> bytes
    pub bytes_inodes_ratio: u32,

    // FS block size
    pub block_size: u32,

    // Inode size
    pub inode_size: u32,
}

impl Default for Params {
    fn default() -> Self {
        Params {
            disk_size: 0,
            bytes_inodes_ratio: 16384,
            block_size: 4096,
            inode_size: 64,
        }
    }
}

/// Superblock contains FS metadata
#[derive(Serialize, Deserialize, PartialEq, Debug)]
pub struct Superblock {
    // Contains MAGIC value
    pub magic: u64,
    // FS block size
    pub block_size: u32,
    // Total number of inodes
    pub num_inodes: u32,
    // Size of the inode
    pub inode_size: u32,
    // Number of inode bitmap blocks
    pub num_inode_bitmap_blocks: u32,
    // Number of data bitmap blocks
    pub num_data_bitmap_blocks: u32,
    // Number of inode blocks
    pub num_inode_blocks: u32,
    // The offset of the first block containing data in bytes
    pub first_data_block_offset: u64,
}

impl Superblock {
    /// Superblock magic number
    pub const MAGIC: u64 = 0x0073666e7261656c;

    /// Create a brand new superblock instance
    pub fn new(params: Params) -> Self {
        let block_bits = params.block_size as f64 * 8.;

        // Calculate the number of inodes
        // For every `params.bytes_inode_ratio` bytes of the disk capacity
        // we need to provision one inode
        let num_inodes = (params.disk_size as f64 /
            params.bytes_inodes_ratio as f64).floor() as u32;

        // How many bitmap blocks we need to track all the required inodes
        // One bitmap block can track `block_size * 8` inodes
        let num_inode_bitmap_blocks = (num_inodes as f64 /
            block_bits).ceil() as u32;

        // How many blocks required for inodes themselves
        let num_inode_blocks = ((params.inode_size * num_inodes) as f64 /
            params.block_size as f64).ceil() as u32;

        // The number of blocks occupied by non-data stuff
        let non_data_blocks = (
            (BOOT_BLOCK_SIZE + SUPERBLOCK_SIZE) as f32 /
                params.block_size as f32).ceil() as u32 +
            num_inode_bitmap_blocks +
            num_inode_blocks;

        // This is the total number of blocks disk has
        let total_blocks = params.disk_size / params.block_size as u64;

        // How many actual data blocks are left
        let num_data_blocks = total_blocks - non_data_blocks as u64;

        // How many data bitmap blocks are needed to map data blocks
        let num_data_bitmap_blocks = (num_data_blocks as f64 /
            block_bits).ceil() as u32;

        let first_data_block_offset = BOOT_BLOCK_SIZE +
            SUPERBLOCK_SIZE +
            (num_inode_bitmap_blocks * params.block_size) as u64 +
            (num_data_bitmap_blocks * params.block_size) as u64 +
            (num_inode_blocks * params.block_size) as u64;

        Superblock {
            magic: Superblock::MAGIC,
            block_size: params.block_size,
            inode_size: params.inode_size,
            num_inodes,
            num_inode_bitmap_blocks,
            num_data_bitmap_blocks,
            num_inode_blocks,
            first_data_block_offset,
        }
    }

    /// Load superblock from raw bytes
    pub fn load(buf: &[u8]) -> Result<Self, Error> {
        let sb: Self = bincode::deserialize(buf)?;

        if sb.magic != Superblock::MAGIC {
            return Err(format_err!("invalid superblock magic: {}", sb.magic));
        }

        return Ok(sb);
    }

    #[inline(always)]
    pub fn block_offset(&self, block: u32) -> u64 {
        self.first_data_block_offset + (block * self.block_size) as u64
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_superblock_create() {
        struct TestCase {
            disk_size: u64,
            block_size: u32,
            inode_size: u32,
            bytes_inodes_ratio: u32,
            num_inodes: u32,
            num_inode_bitmap_blocks: u32,
            num_data_bitmap_blocks: u32,
            num_inode_blocks: u32,
            first_data_block_offset: u64,
        };

        let cases = vec![
            TestCase {
                block_size: 512u32,
                inode_size: 64u32,
                disk_size: 1234567890,
                bytes_inodes_ratio: 16384,
                // Number of inodes
                // floor(1234567890 / 16384) = floor(75352.044) = 75352
                num_inodes: 75352,
                // Number of inode bitmap blocks
                // ceil(75352 / (512 * 8)) = ceil(18.4) = 19
                num_inode_bitmap_blocks: 19,
                // Number of data bitmap blocks
                // ceil(2401824 / (512*8)) = ceil(586.38) = 587
                num_data_bitmap_blocks: 587,
                // Number of inode blocks
                // ceil(64 * 75352 / 512) = ceil(9419.0) = 9419
                num_inode_blocks: 9419,
                first_data_block_offset: 5134848,
            }
        ];

        for case in cases {
            let params = Params {
                disk_size: case.disk_size,
                bytes_inodes_ratio: case.bytes_inodes_ratio,
                block_size: case.block_size,
                inode_size: case.inode_size,
            };

            let sb = Superblock::new(params);

            assert_eq!(sb.magic, Superblock::MAGIC);
            assert_eq!(sb.block_size, case.block_size);
            assert_eq!(sb.inode_size, case.inode_size);

            // Verify calculated properties
            assert_eq!(sb.num_inodes, case.num_inodes);
            assert_eq!(sb.num_inode_bitmap_blocks,
                       case.num_inode_bitmap_blocks);
            assert_eq!(sb.num_data_bitmap_blocks, case.num_data_bitmap_blocks);
            assert_eq!(sb.num_inode_blocks, case.num_inode_blocks);
            assert_eq!(sb.first_data_block_offset,
                       case.first_data_block_offset);
        }
    }
}
