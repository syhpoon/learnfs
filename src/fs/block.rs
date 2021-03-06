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

use super::types::BlockPtr;

/// Block structure represents a single block of data
pub struct Block {
    pub idx: BlockPtr,
    pub data: Vec<u8>,
    pub dirty: bool,
}

impl Block {
    /// Create an empty block of a given size
    pub fn new(idx: BlockPtr, size: usize) -> Self {
        Block {
            idx,
            data: vec![0u8; size],
            dirty: false,
        }
    }

    /// Create a block from existing data
    pub fn from_buf(idx: BlockPtr, buf: Vec<u8>) -> Self {
        Block {
            idx,
            data: buf,
            dirty: false,
        }
    }

    pub fn size(&self) -> usize {
        self.data.len()
    }

    pub fn set_dirty(&mut self, dirty: bool) {
        self.dirty = dirty
    }

    pub fn is_dirty(&self) -> bool {
        self.dirty
    }
}
