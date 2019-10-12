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

use std::io::SeekFrom;

use failure::Error;

use super::{Bitmap, Superblock};
use crate::device::Device;
use crate::fs;

#[derive(Debug)]
pub struct FsInfo {
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

/// Filesystem structure is the main entry point to all
/// the filesystem operations
pub struct Filesystem<T: Device> {
    device: T,
    superblock: Superblock,
    inode_bitmap: Bitmap,
    data_bitmap: Bitmap,
}

impl<T: Device> Filesystem<T> {
    /// Create a new filesystem on a given device.
    /// Existing data on the device will become corrupted/unavailable.
    pub fn create(mut device: T) -> Result<Self, Error> {
        // Skip boot block
        device.seek(SeekFrom::Start(fs::BOOT_BLOCK_SIZE))?;

        // Write the superblock first
        let sb = Filesystem::write_superblock(&mut device)?;

        // Write inode bitmap
        let inode_bitmap = Bitmap::new(sb.inode_bitmap_size());
        device.write_all(inode_bitmap.as_slice())?;

        // Write data bitmap
        let data_bitmap = Bitmap::new(sb.data_bitmap_size());
        device.write_all(data_bitmap.as_slice())?;

        // Next create a root directory
        // TODO
        // 1. Get the next free inode
        // 2. Allocate a new block
        // 3. Add two dir entries: . and ..

        let fs = Filesystem {
            device,
            inode_bitmap,
            data_bitmap,
            superblock: sb,
        };

        Ok(fs)
    }

    /// Load existing filesystem from a device
    pub fn load(mut device: T) -> Result<Self, Error> {
        // Skip boot block
        device.seek(SeekFrom::Start(fs::BOOT_BLOCK_SIZE))?;

        let mut buf: Vec<u8> = vec![0; fs::SUPERBLOCK_SIZE as usize];
        device.read(buf.as_mut_slice())?;

        let sb = Superblock::load(buf.as_slice())?;

        // Load bitmaps
        let inode_bitmap = Filesystem::load_bitmap(sb.inode_bitmap_size(),
                                                   &mut device)?;

        let data_bitmap = Filesystem::load_bitmap(sb.data_bitmap_size(),
                                                  &mut device)?;

        let fs = Filesystem {
            device,
            inode_bitmap,
            data_bitmap,
            superblock: sb,
        };

        Ok(fs)
    }

    pub fn info(&self) -> FsInfo {
        FsInfo {
            block_size: self.superblock.block_size,
            num_inodes: self.superblock.num_inodes,
            inode_size: self.superblock.inode_size,
            num_inode_bitmap_blocks: self.superblock.num_inode_bitmap_blocks,
            num_data_bitmap_blocks: self.superblock.num_data_bitmap_blocks,
            num_inode_blocks: self.superblock.num_inode_blocks,
            first_data_block_offset: self.superblock.first_data_block_offset,
        }
    }

    fn write_superblock(device: &mut T) -> Result<Superblock, Error> {
        let sb_params = fs::superblock::Params {
            disk_size: device.capacity()?,
            ..Default::default()
        };

        let sb = fs::Superblock::new(sb_params);

        // Write superblock
        let bytes: Vec<u8> = bincode::serialize(&sb)?;
        device.write(bytes.as_slice())?;

        Ok(sb)
    }

    fn load_bitmap(size: usize, device: &mut T) -> Result<Bitmap, Error> {
        let mut buf: Vec<u8> = vec![0; size];

        device.read(buf.as_mut_slice())?;

        Ok(Bitmap::from_buf(buf))
    }
}
