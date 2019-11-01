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
    // The offset of the first block containing inode in bytes
    pub first_inode_block_offset: u64,
    // The offset of the first block containing data in bytes
    pub first_data_block_offset: u64,
}

