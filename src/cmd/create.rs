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

use std::io;
use std::io::{Write, Seek, SeekFrom};

use crate::fs;

/// Create a new filesystem on a given device
pub fn create_fs<T>(device: &mut T) -> io::Result<()>
    where T: Write + Seek {

    create_superblock(device)?;

    Ok(())
}

fn create_superblock<T>(device: &mut T) -> io::Result<()>
    where T: Write + Seek {

    let sb_params: fs::superblock::Params = Default::default();
    let sb = fs::Superblock::new(sb_params);

    // Skip boot block
    device.seek(SeekFrom::Start(fs::BOOT_BLOCK_SIZE)).unwrap();

    // Write superblock
    let bytes: Vec<u8> = bincode::serialize(&sb).unwrap();
    device.write(bytes.as_slice()).unwrap();

    // Zero inode and data bitmaps
    let bitmap_size = sb.num_inode_bitmap_blocks * sb.block_size;

    // Allocate memory for inode bitmap blocks
    let mut inode_bitmap: Vec<u8> = vec![0; bitmap_size as usize];

    // Mark zeroth inode block as used because block numbering stars with 1
    inode_bitmap[0] = 1;

    // The first inode is used for the root dir
    inode_bitmap[1] = 1;

    device.write_all(inode_bitmap.as_slice()).unwrap();

    Ok(())
}
