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

use failure::{Error, format_err};
use serde::{Deserialize, Serialize};

use super::types::InodePtr;

pub const MAX_FILE_NAME_SIZE: usize = 60;

big_array!{
   BigArray;
   60,
}

#[derive(Serialize, Deserialize)]
pub struct DirEntry {
    // Entry inode number
    pub inode: u32,

    // File name
    #[serde(with = "BigArray")]
    pub name: [u8; MAX_FILE_NAME_SIZE],
}

impl DirEntry {
    pub const SIZE: usize = 64;

    pub fn new(inode: InodePtr, name: &str) -> Result<Self, Error> {
        let len = name.len();
        if len > MAX_FILE_NAME_SIZE {
            return Err(format_err!("file name too long"));
        }

        let mut arr: [u8; MAX_FILE_NAME_SIZE] = [0u8; MAX_FILE_NAME_SIZE];
        (&mut arr[..len]).copy_from_slice(name.as_bytes());

        Ok(DirEntry { inode, name: arr })
    }

    pub fn name_as_string(&self) -> String {
        // TODO: do not clone but return slice instead
        return String::from_utf8(
            self.name.iter().take_while(|&ch| *ch != 0)
                .cloned()
                .collect()
        ).unwrap();
    }
}

impl TryFrom<&[u8]> for DirEntry {
    type Error = Error;

    fn try_from(data: &[u8]) -> Result<Self, Self::Error> {
        let entry: Self = bincode::deserialize(data)?;

        Ok(entry)
    }
}
