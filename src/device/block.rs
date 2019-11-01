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

use std::fs;
use std::io;

use failure::Error;
use std::os::unix::io::AsRawFd;

use crate::device::Device;
use std::fs::{OpenOptions, File};

// Generate ioctl function
const BLKGETSIZE64_CODE: u8 = 0x12; // Defined in linux/fs.h
const BLKGETSIZE64_SEQ: u8 = 114;
ioctl_read!(ioctl_blkgetsize64, BLKGETSIZE64_CODE, BLKGETSIZE64_SEQ, u64);

/// Implements Device abstraction on top of a block device
pub struct BlockDevice {
    file: fs::File,
    capacity: u64,
}

impl BlockDevice {
    /// Instantiate new BlockDevice from a path
    pub fn new(path: &str) -> Result<Self, Error> {
        let file = OpenOptions::new()
            .write(true)
            .open(path)?;

        let fd = file.as_raw_fd();

        // Determine device capacity
        let mut cap = 0u64;
        let cap_ptr = &mut cap as *mut u64;

        unsafe {
            ioctl_blkgetsize64(fd, cap_ptr)?;
        };

        Ok(BlockDevice { file, capacity: cap })
    }

    /// Open existing file
    pub fn open(path: &str) -> Result<Self, Error> {
        let file = File::open(path)?;

        Ok(BlockDevice { file, capacity: 0 })
    }
}

impl Device for BlockDevice {
    fn capacity(&mut self) -> Result<u64, Error> {
        Ok(self.capacity)
    }
}

impl io::Write for BlockDevice {
    fn write(&mut self, buf: &[u8]) -> Result<usize, io::Error> {
        return self.file.write(buf);
    }

    fn flush(&mut self) -> Result<(), io::Error> {
        return self.file.flush();
    }
}

impl io::Seek for BlockDevice {
    fn seek(&mut self, pos: io::SeekFrom) -> Result<u64, io::Error> {
        return self.file.seek(pos);
    }
}

impl io::Read for BlockDevice {
    fn read(&mut self, buf: &mut [u8]) -> Result<usize, io::Error> {
        return self.file.read(buf);
    }
}
