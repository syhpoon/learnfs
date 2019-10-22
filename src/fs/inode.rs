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

use nix::sys::stat::{mode_t, SFlag, Mode};
use lazy_static;

pub const DIRECT_BLOCKS_NUM: usize = 21;

lazy_static! {
   static ref S_IFMT: mode_t = SFlag::S_IFMT.bits();
   static ref S_IFDIR: mode_t = SFlag::S_IFDIR.bits();
   static ref S_IFREG: mode_t = SFlag::S_IFREG.bits();
   static ref S_IFLNK: mode_t = SFlag::S_IFLNK.bits();
   static ref S_IFSOCK: mode_t = SFlag::S_IFSOCK.bits();
   static ref S_IFIFO: mode_t = SFlag::S_IFIFO.bits();
   static ref S_IFBLK: mode_t = SFlag::S_IFBLK.bits();
   static ref S_IFCHR: mode_t = SFlag::S_IFCHR.bits();

   static ref S_ISUID: mode_t = Mode::S_ISUID.bits();
   static ref S_ISGID: mode_t = Mode::S_ISGID.bits();
   static ref S_ISVTX: mode_t = Mode::S_ISVTX.bits();
   static ref S_IRWXU: mode_t = Mode::S_IRWXU.bits();
   static ref S_IRUSR: mode_t = Mode::S_IRUSR.bits();
   static ref S_IWUSR: mode_t = Mode::S_IWUSR.bits();
   static ref S_IXUSR: mode_t = Mode::S_IXUSR.bits();
   static ref S_IRWXG: mode_t = Mode::S_IRWXG.bits();
   static ref S_IRGRP: mode_t = Mode::S_IRGRP.bits();
   static ref S_IWGRP: mode_t = Mode::S_IWGRP.bits();
   static ref S_IXGRP: mode_t = Mode::S_IXGRP.bits();
   static ref S_IRWXO: mode_t = Mode::S_IRWXO.bits();
   static ref S_IROTH: mode_t = Mode::S_IROTH.bits();
   static ref S_IWOTH: mode_t = Mode::S_IWOTH.bits();
   static ref S_IXOTH: mode_t = Mode::S_IXOTH.bits();
}

#[derive(PartialEq, Debug, Copy, Clone)]
pub struct InodePerm(pub mode_t);

impl InodePerm {
    #[inline(always)]
    pub fn set_user_id() -> mode_t {
        return *S_ISUID;
    }

    #[inline(always)]
    pub fn set_group_id() -> mode_t {
        return *S_ISGID;
    }

    #[inline(always)]
    pub fn sticky_bit() -> mode_t {
        return *S_ISVTX;
    }

    #[inline(always)]
    pub fn owner_read_write_execute() -> mode_t {
        return *S_IRWXU;
    }

    #[inline(always)]
    pub fn owner_read() -> mode_t {
        return *S_IRUSR;
    }

    #[inline(always)]
    pub fn owner_write() -> mode_t {
        return *S_IWUSR;
    }

    #[inline(always)]
    pub fn owner_execute() -> mode_t {
        return *S_IXUSR;
    }

    #[inline(always)]
    pub fn group_read_write_execute() -> mode_t {
        return *S_IRWXG;
    }

    #[inline(always)]
    pub fn group_read() -> mode_t {
        return *S_IRGRP;
    }

    #[inline(always)]
    pub fn group_write() -> mode_t {
        return *S_IWGRP;
    }

    #[inline(always)]
    pub fn group_execute() -> mode_t {
        return *S_IXGRP;
    }

    #[inline(always)]
    pub fn others_read_write_execute() -> mode_t {
        return *S_IRWXO;
    }

    #[inline(always)]
    pub fn others_read() -> mode_t {
        return *S_IROTH;
    }

    #[inline(always)]
    pub fn others_write() -> mode_t {
        return *S_IWOTH;
    }

    #[inline(always)]
    pub fn others_execute() -> mode_t {
        return *S_IXOTH;
    }
}

#[derive(PartialEq, Debug, Copy, Clone)]
pub enum InodeType {
    Directory,
    File,
    Link,
    Socket,
    Fifo,
    Block,
    Char,
}

impl InodeType {
    pub fn bits(&self) -> mode_t {
        use InodeType::*;

        match self {
            Directory => *S_IFDIR,
            File => *S_IFREG,
            Link => *S_IFLNK,
            Socket => *S_IFSOCK,
            Fifo => *S_IFIFO,
            Block => *S_IFBLK,
            Char => *S_IFCHR,
        }
    }

    pub fn from_bits(bits: mode_t) -> Option<Self> {
        use InodeType::*;

        match bits {
            v if v == *S_IFDIR => Some(Directory),
            v if v == *S_IFREG => Some(File),
            v if v == *S_IFLNK => Some(Link),
            v if v == *S_IFSOCK => Some(Socket),
            v if v == *S_IFIFO => Some(Fifo),
            v if v == *S_IFBLK => Some(Block),
            v if v == *S_IFCHR => Some(Char),
            _ => None
        }
    }
}

#[derive(Debug)]
pub struct Inode {
    // File type and mode
    pub mode: mode_t,
    // Inode number
    pub ino: u32,
    // Number of directories containing this inode
    pub n_link: u32,
    // Owner's user id
    pub uid: u32,
    // Owner's group id
    pub gid: u32,
    // File size in bytes
    pub size: u32,
    // Access time
    pub atime: u32,
    // Modification time
    pub mtime: u32,
    // Status change time
    pub ctime: u32,
    // Direct block pointers
    blocks: [u32; DIRECT_BLOCKS_NUM],
    // Single indirect block pointers
    indirect_block: u32,
    // Double indirect block pointers
    double_indirect_block: u32,
}

impl Inode {
    /// Create a new empty inode
    pub fn new(ino: u32) -> Self {
        Inode {
            ino,
            mode: 0,
            n_link: 0,
            uid: 0,
            gid: 0,
            size: 0,
            atime: 0,
            mtime: 0,
            ctime: 0,
            blocks: [0u32; DIRECT_BLOCKS_NUM],
            indirect_block: 0,
            double_indirect_block: 0,
        }
    }

    /// Set inode type
    pub fn set_type(&mut self, t: InodeType) {
        self.mode = (self.mode & !*S_IFMT) | (t.bits() & *S_IFMT)
    }

    /// Return inode type
    pub fn get_type(&self) -> Option<InodeType> {
        InodeType::from_bits(self.mode & *S_IFMT)
    }

    /// Set inode permission
    pub fn set_perm(&mut self, p: InodePerm) {
        self.mode = (self.mode & *S_IFMT) | (p.0 & !*S_IFMT)
    }

    /// Return inode permissions
    pub fn get_perm(&self) -> InodePerm {
        InodePerm(self.mode & !*S_IFMT)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use InodeType::*;

    #[test]
    fn test_inode_type() {
        let mut inode = Inode::new(0);

        assert_eq!(inode.get_type(), None);

        let types = vec![
            Directory,
            File,
            Link,
            Socket,
            Fifo,
            Block,
            Char,
        ];

        for t in types {
            inode.set_type(t);
            assert_eq!(inode.get_type(), Some(t));
        }
    }

    #[test]
    fn test_inode_perm() {
        let mut inode = Inode::new(0);

        let perms = vec![
            InodePerm::set_user_id(),
            InodePerm::set_group_id(),
            InodePerm::sticky_bit(),
            InodePerm::owner_read_write_execute(),
            InodePerm::owner_read(),
            InodePerm::owner_write(),
            InodePerm::owner_execute(),
            InodePerm::group_read_write_execute(),
            InodePerm::group_read(),
            InodePerm::group_write(),
            InodePerm::group_execute(),
            InodePerm::others_read_write_execute(),
            InodePerm::others_read(),
            InodePerm::others_write(),
            InodePerm::others_execute(),
            InodePerm::owner_read_write_execute() |
                InodePerm::group_read() | InodePerm::group_execute() |
                InodePerm::others_write(),
        ];

        for p in perms {
            let p = InodePerm(p);

            inode.set_perm(p);
            assert_eq!(inode.get_perm(), p);
        }
    }

    #[test]
    fn test_inode_type_and_mode() {
        let mut inode = Inode::new(0);

        assert_eq!(inode.get_type(), None);

        let types = vec![
            (Directory, InodePerm::owner_read_write_execute() |
                InodePerm::sticky_bit()),
            (File, InodePerm::owner_read() |
                InodePerm::group_read_write_execute() |
                InodePerm::others_read()),
            (Block, InodePerm::owner_read_write_execute() |
                InodePerm::group_read_write_execute() |
                InodePerm::group_read_write_execute()),
        ];

        for (t, p) in types {
            let p = InodePerm(p);
            inode.set_type(t);
            inode.set_perm(p);

            assert_eq!(inode.get_type(), Some(t));
            assert_eq!(inode.get_perm(), p);
        }
    }
}
