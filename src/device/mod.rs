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

mod block;
mod file;

use std::io;

use failure::{Error, format_err};
use nix::sys::stat::{stat, SFlag};

pub use file::FileDevice;
pub use block::BlockDevice;

/// Defines a device abstraction.
/// All the rest of the code works with this interface.
pub trait Device: io::Write + io::Read + io::Seek {
    /// Return device capacity in bytes
    fn capacity(&mut self) -> Result<u64, Error>;
}

pub enum DeviceType {
    File,
    Block,
}

pub fn device_type(path: &str) -> Result<DeviceType, Error> {
    match stat(path) {
        Ok(st) => {
            match SFlag::from_bits(st.st_mode & SFlag::S_IFMT.bits()) {
                // Block device
                Some(SFlag::S_IFBLK) => Ok(DeviceType::Block),
                // Regular file
                Some(SFlag::S_IFREG) => Ok(DeviceType::File),
                _ => Err(format_err!(
                    "Only block devices and regular files are currently supported"
                )),
            }
        }
        Err(e) => Err(format_err!("Failed to read device stat: {}", e)),
    }
}

