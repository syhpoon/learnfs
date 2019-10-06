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

use crate::device::Device;
use std::fs::{OpenOptions, File};

/// Implements Device abstraction on top of an existing file
pub struct FileDevice {
    file: fs::File,
    capacity: u64,
}

impl FileDevice {
    /// Instantiate new FileDevice from a file path
    pub fn new(path: &str, capacity: u64) -> Result<Self, Error> {
        let mut file = OpenOptions::new()
            .write(true)
            .truncate(true)
            .create(true)
            .open(path)?;

        Ok(FileDevice { file, capacity })
    }

    /// Open existing files
    pub fn open(path: &str) -> Result<Self, Error> {
        let file = File::open(path)?;

        Ok(FileDevice { file, capacity: 0 })
    }
}

impl Device for FileDevice {
    fn capacity(&mut self) -> Result<u64, Error> {
        Ok(self.capacity)
    }
}

impl io::Write for FileDevice {
    fn write(&mut self, buf: &[u8]) -> Result<usize, io::Error> {
        return self.file.write(buf);
    }

    fn flush(&mut self) -> Result<(), io::Error> {
        return self.file.flush();
    }
}

impl io::Seek for FileDevice {
    fn seek(&mut self, pos: io::SeekFrom) -> Result<u64, io::Error> {
        return self.file.seek(pos);
    }
}

impl io::Read for FileDevice {
    fn read(&mut self, buf: &mut [u8]) -> Result<usize, io::Error> {
        return self.file.read(buf);
    }
}

