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

use std::io::{SeekFrom, Read};

use failure::{Error, format_err};

use crate::device::Device;
use crate::fs;
use crate::fs::Superblock;

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

/// Return filesystem info
pub fn fs_info(mut device: impl Device) -> Result<FsInfo, Error> {
    // Skip boot block
    device.seek(SeekFrom::Start(fs::BOOT_BLOCK_SIZE))?;

    let mut buf: Vec<u8> = vec![0; fs::SUPERBLOCK_SIZE as usize];
    device.read(buf.as_mut_slice())?;

    let sb: fs::Superblock = bincode::deserialize(buf.as_slice())?;

    if sb.magic != Superblock::MAGIC {
        return Err(format_err!("invalid superblock magic: {}", sb.magic));
    }

    Ok(FsInfo {
        block_size: sb.block_size,
        num_inodes: sb.num_inodes,
        inode_size: sb.inode_size,
        num_inode_bitmap_blocks: sb.num_inode_bitmap_blocks,
        num_data_bitmap_blocks: sb.num_data_bitmap_blocks,
        num_inode_blocks: sb.num_inode_blocks,
        first_data_block_offset: sb.first_data_block_offset,
    })
}
