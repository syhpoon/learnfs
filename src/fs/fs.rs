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
use std::collections::{HashMap, VecDeque};
use std::rc::Rc;
use std::cell::RefCell;

use failure::{Error, format_err};

use crate::fs;
use crate::device::Device;
use super::types::*;
use super::{Bitmap, Block, Directory, Inode, Superblock};
use super::inode::{InodeType, InodePerm};
use super::fs_info::FsInfo;

/// Filesystem structure is the main entry point to all
/// the filesystem operations
pub struct Filesystem<T: Device> {
    device: T,
    superblock: Superblock,
    inode_bitmap: Bitmap,
    data_bitmap: Bitmap,
    next_free_inode: Option<usize>,
    next_free_block: Option<usize>,
    block_cache: HashMap<BlockPtr, Rc<RefCell<Block>>>,
    inode_cache: HashMap<InodePtr, Rc<RefCell<Inode>>>,
    //dir_cache: HashMap<InodePtr, Rc<RefCell<Inode>>>,
}

impl<T: Device> Filesystem<T> {
    /// Create a new filesystem on a given device.
    /// Existing data on the device will be destroyed.
    pub fn create(mut device: T) -> Result<Self, Error> {
        // Skip boot block
        device.seek(SeekFrom::Start(fs::BOOT_BLOCK_SIZE))?;

        // Create superblock
        let superblock = Filesystem::create_superblock(&mut device)?;

        // Write inode bitmap
        let mut inode_bitmap = Bitmap::new(superblock.inode_bitmap_size());

        // Zeroth inode is not available
        inode_bitmap.set(0);
        device.write_all(inode_bitmap.as_slice())?;

        // Write data bitmap
        let mut data_bitmap = Bitmap::new(superblock.data_bitmap_size());

        // Zeroth block is not available
        data_bitmap.set(0);
        device.write_all(data_bitmap.as_slice())?;

        let mut fs = Filesystem {
            device,
            inode_bitmap,
            data_bitmap,
            superblock,
            next_free_inode: Some(1),
            next_free_block: Some(1),
            inode_cache: HashMap::new(),
            block_cache: HashMap::new(),
        };

        // We need to have a root directory on a newly created filesystem
        let mut inode = fs.allocate_inode()?;
        inode.set_dirty(true);

        inode.set_type(InodeType::Directory);
        inode.set_perm(
            InodePerm(
                InodePerm::owner_read_write_execute() |
                    InodePerm::group_read() |
                    InodePerm::group_execute() |
                    InodePerm::others_read() |
                    InodePerm::others_execute()));

        let block = fs.allocate_block()?;

        // The inode now has one data block
        inode.add_block(block.idx)?;

        // Put block into cache
        fs.put_block(Rc::new(RefCell::new(block)));

        let mut dir = Directory::load(&inode, &mut fs)?;

        dir.add_entry(String::from("."), inode.ino, &mut fs)?;
        dir.add_entry(String::from(".."), inode.ino, &mut fs)?;

        fs.put_inode(Rc::new(RefCell::new(inode)));

        //TODO: Add dir to cache?

        // Flush changes to device
        fs.flush()?;

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
        let mut inode_bitmap = Filesystem::load_bitmap(sb.inode_bitmap_size(),
                                                       &mut device)?;

        let mut data_bitmap = Filesystem::load_bitmap(sb.data_bitmap_size(),
                                                      &mut device)?;

        let next_free_inode = inode_bitmap.next_clear_bit(0);
        let next_free_block = data_bitmap.next_clear_bit(0);

        let fs = Filesystem {
            device,
            inode_bitmap,
            data_bitmap,
            next_free_inode,
            next_free_block,
            superblock: sb,
            block_cache: HashMap::new(),
            inode_cache: HashMap::new(),
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

    fn create_superblock(device: &mut T) -> Result<Superblock, Error> {
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

    fn allocate_inode(&mut self) -> Result<Inode, Error> {
        let idx = match self.next_free_inode {
            Some(idx) => idx,
            None => return Err(format_err!("no free inodes left"))
        };

        // Mark inode as used
        self.inode_bitmap.set(idx);

        let inode = Inode::new(idx as u32);
        self.next_free_inode = self.inode_bitmap.next_clear_bit(idx);

        Ok(inode)
    }

    fn allocate_block(&mut self) -> Result<Block, Error> {
        let idx = match self.next_free_block {
            Some(idx) => idx,
            None => return Err(format_err!("no free blocks left"))
        };

        // Mark block as used
        self.data_bitmap.set(idx);

        let block = Block::new(idx as u32,
                               self.superblock.block_size as usize);
        self.next_free_block = self.data_bitmap.next_clear_bit(idx);

        Ok(block)
    }

    #[inline(always)]
    fn block_offset(&self, block_idx: BlockPtr) -> u64 {
        self.superblock.first_data_block_offset +
            (block_idx * self.superblock.block_size) as u64
    }

    #[inline(always)]
    fn inode_offset(&self, inode_idx: InodePtr) -> u64 {
        self.superblock.first_inode_block_offset +
            (inode_idx * self.superblock.inode_size) as u64
    }

    /// Put block into a cache
    fn put_block(&mut self, block: Rc<RefCell<Block>>) {
        let idx = block.borrow().idx;
        self.block_cache.insert(idx, block);
    }

    /// Put inode into a cache
    fn put_inode(&mut self, inode: Rc<RefCell<Inode>>) {
        let idx = inode.borrow().ino;
        self.inode_cache.insert(idx, inode);
    }

    /// Return a given block.
    /// The block will either be returned from a cache or loaded from device
    pub fn get_block(&mut self, idx: BlockPtr) -> Result<
        Rc<RefCell<Block>>, Error> {
        match self.block_cache.get_mut(&idx) {
            Some(block) => Ok(Rc::clone(&block)),
            None => {
                println!("get_block ---------- idx={}", idx);
                // Load from device
                let mut buf: Vec<u8> = vec![
                    0; self.superblock.block_size as usize];

                self.device.seek(SeekFrom::Start(self.block_offset(idx)))?;
                self.device.read(buf.as_mut_slice())?;

                let block = Block::from_buf(idx, buf);

                let brc = Rc::new(RefCell::new(block));
                let res = Rc::clone(&brc);

                self.put_block(brc);

                Ok(res)
            }
        }
    }

    fn flush(&mut self) -> Result<(), Error> {
        let mut blocks: VecDeque<Rc<RefCell<Block>>> = VecDeque::new();
        let mut inodes: VecDeque<Rc<RefCell<Inode>>> = VecDeque::new();

        for cell in self.block_cache.values() {
            if cell.borrow().is_dirty() {
                blocks.push_back(Rc::clone(&cell));
            }
        }

        for cell in self.inode_cache.values() {
            if cell.borrow().is_dirty() {
                inodes.push_back(Rc::clone(&cell));
            }

            //TODO: inode blocks?
        }

        for cell in blocks {
            let mut block = cell.borrow_mut();
            self.flush_block(&block)?;
            block.set_dirty(false);
        };

        for cell in inodes {
            let mut inode = cell.borrow_mut();
            self.flush_inode(&inode)?;
            inode.set_dirty(false);
        };

        Ok(())
    }

    fn flush_block(&mut self, block: &Block) -> Result<(), Error> {
        self.device.seek(SeekFrom::Start(self.block_offset(block.idx)))?;
        self.device.write(block.data.as_slice())?;

        Ok(())
    }

    fn flush_inode(&mut self, inode: &Inode) -> Result<(), Error> {
        let bytes: Vec<u8> = bincode::serialize(&inode)?;
        self.device.seek(SeekFrom::Start(self.inode_offset(inode.ino)))?;

        (&mut self.device).write(bytes.as_slice())?;

        Ok(())
    }
}
