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

use std::convert::TryFrom;
use std::collections::HashMap;

use failure::Error;

use crate::device::Device;
use super::{Filesystem, Inode};
use super::types::InodePtr;
use super::dir_entry::DirEntry;
use super::types::BlockPtr;

pub struct Directory {
    name2inode: HashMap<String, InodePtr>,
    entries: Vec<(Option<DirEntry>, BlockPtr)>,
}

impl Directory {
    /// Load directory entries
    pub fn load<T: Device>(inode: &Inode,
                           fs: &mut Filesystem<T>) -> Result<Self, Error> {
        let mut dir = Directory {
            name2inode: HashMap::new(),
            entries: vec![],
        };

        // Load every assigned block and scan for dir entries
        for bid in inode.get_blocks() {
            let bcell = fs.get_block(bid)?;
            let block = bcell.borrow();

            let mut start = 0;
            let mut end = DirEntry::SIZE;

            while end < block.size() {
                let slice = &block.data[start..end];
                let entry = DirEntry::try_from(slice)?;

                // inode == 0 indicates unused entry
                if entry.inode == 0 {
                    dir.entries.push((None, bid));
                } else {
                    dir.name2inode.insert(entry.name_as_string(), entry.inode);
                    dir.entries.push((Some(entry), bid));
                }

                start += DirEntry::SIZE;
                end += DirEntry::SIZE;
            }
        }

        Ok(dir)
    }

    pub fn add_entry<T: Device>(&mut self, name: String,
                                inode: InodePtr,
                                fs: &mut Filesystem<T>) -> Result<(), Error> {
        let entry = DirEntry::new(inode, name.as_str())?;
        self.name2inode.insert(name, inode);

        // Find the first empty entry slot
        match self.entries.iter().position(|e| e.0.is_none()) {
            Some(idx) => {
                let bid = self.entries[idx].1;
                let offset = idx * DirEntry::SIZE;
                let bytes: Vec<u8> = bincode::serialize(&entry)?;
                let range = offset..(offset + DirEntry::SIZE);
                let bcell = fs.get_block(bid)?;
                let mut block = bcell.borrow_mut();

                // Replace serialized entry
                block.data.splice(range, bytes);
                block.set_dirty(true);

                self.entries[idx] = (Some(entry), bid)
            }
            None => unimplemented!("TODO: allocate more blocks")
        };

        Ok(())
    }
}
