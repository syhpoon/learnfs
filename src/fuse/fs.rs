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

use std::ffi::OsStr;
use std::path::Path;
use time::Timespec;

use failure::Error;
use log;
use fuse as f;
use libc::ENOSYS;
use tonic::transport::channel::Channel;

use crate::grpc::Client;

type Cl = Client<Channel>;
pub type Session = fuse::BackgroundSession<'static>;

pub struct FuseFs {
    client: Cl,
}

impl f::Filesystem for FuseFs {
    fn destroy(&mut self, _req: &f::Request) {}

    /// Look up a directory entry by name and get its attributes.
    fn lookup(&mut self,
              _req: &f::Request,
              _parent: u64,
              _name: &OsStr,
              reply: f::ReplyEntry) {
        log::info!("@@@ lookup");
        log::info!("    req={:?}", _req);
        log::info!("    parent={:?}", _parent);
        log::info!("    _name={:?}", _name);
        reply.error(ENOSYS);
    }

    /// Forget about an inode.
    /// The nlookup parameter indicates the number of lookups previously performed on
    /// this inode. If the filesystem implements inode lifetimes, it is recommended that
    /// inodes acquire a single reference on each lookup, and lose nlookup references on
    /// each forget. The filesystem may ignore forget calls, if the inodes don't need to
    /// have a limited lifetime. On unmount it is not guaranteed, that all referenced
    /// inodes will receive a forget message.
    fn forget(&mut self, _req: &f::Request, _ino: u64, _nlookup: u64) {}

    /// Get file attributes.
    fn getattr(&mut self, _req: &f::Request, _ino: u64, reply: f::ReplyAttr) {
        log::info!("@@@ getattr");
        log::info!("    req={:?}", _req);
        log::info!("    ino={:?}", _ino);
        reply.error(ENOSYS);
    }

    /// Set file attributes.
    fn setattr(&mut self,
               _req: &f::Request,
               _ino: u64,
               _mode: Option<u32>,
               _uid: Option<u32>,
               _gid: Option<u32>,
               _size: Option<u64>,
               _atime: Option<Timespec>,
               _mtime: Option<Timespec>,
               _fh: Option<u64>,
               _crtime: Option<Timespec>,
               _chgtime: Option<Timespec>,
               _bkuptime: Option<Timespec>,
               _flags: Option<u32>,
               reply: f::ReplyAttr) {
        log::info!("@@@ setattr");
        reply.error(ENOSYS);
    }

    /// Read symbolic link.
    fn readlink(&mut self, _req: &f::Request, _ino: u64, reply: f::ReplyData) {
        log::info!("@@@ readlink");
        reply.error(ENOSYS);
    }

    /// Create file node.
    fn mknod(&mut self,
             _req: &f::Request,
             _parent: u64,
             _name: &OsStr,
             _mode: u32,
             _rdev: u32,
             reply: f::ReplyEntry) {
        log::info!("@@@ mknod");
        reply.error(ENOSYS);
    }

    /// Create a directory.
    fn mkdir(&mut self,
             _req: &f::Request,
             _parent: u64,
             _name: &OsStr,
             _mode: u32,
             reply: f::ReplyEntry) {
        log::info!("@@@ mkdir");
        reply.error(ENOSYS);
    }

    /// Remove a file.
    fn unlink(&mut self,
              _req: &f::Request,
              _parent: u64,
              _name: &OsStr,
              reply: f::ReplyEmpty) {
        log::info!("@@@ unlink");
        reply.error(ENOSYS);
    }

    /// Remove a directory.
    fn rmdir(&mut self,
             _req: &f::Request,
             _parent: u64,
             _name: &OsStr,
             reply: f::ReplyEmpty) {
        log::info!("@@@ rmdir");
        reply.error(ENOSYS);
    }

    /// Create a symbolic link.
    fn symlink(&mut self,
               _req: &f::Request,
               _parent: u64,
               _name: &OsStr,
               _link: &Path,
               reply: f::ReplyEntry) {
        log::info!("@@@ symlink");
        reply.error(ENOSYS);
    }

    /// Rename a file.
    fn rename(&mut self,
              _req: &f::Request,
              _parent: u64,
              _name: &OsStr,
              _newparent: u64,
              _newname: &OsStr,
              reply: f::ReplyEmpty) {
        log::info!("@@@ rename");
        reply.error(ENOSYS);
    }

    /// Create a hard link.
    fn link(&mut self,
            _req: &f::Request,
            _ino: u64,
            _newparent: u64,
            _newname: &OsStr,
            reply: f::ReplyEntry) {
        log::info!("@@@ link");
        reply.error(ENOSYS);
    }

    /// Open a file.
    /// Open flags (with the exception of O_CREAT, O_EXCL, O_NOCTTY and O_TRUNC) are
    /// available in flags. Filesystem may store an arbitrary file handle (pointer, index,
    /// etc) in fh, and use this in other all other file operations (read, write, flush,
    /// release, fsync). Filesystem may also implement stateless file I/O and not store
    /// anything in fh. There are also some flags (direct_io, keep_cache) which the
    /// filesystem may set, to change the way the file is opened. See fuse_file_info
    /// structure in <fuse_common.h> for more details.
    fn open(&mut self,
            _req: &f::Request,
            _ino: u64,
            _flags: u32,
            reply: f::ReplyOpen) {
        log::info!("@@@ open");
        reply.opened(0, 0);
    }

    /// Read data.
    /// Read should send exactly the number of bytes requested except on EOF or error,
    /// otherwise the rest of the data will be substituted with zeroes. An exception to
    /// this is when the file has been opened in 'direct_io' mode, in which case the
    /// return value of the read system call will reflect the return value of this
    /// operation. fh will contain the value set by the open method, or will be undefined
    /// if the open method didn't set any value.
    fn read(&mut self,
            _req: &f::Request,
            _ino: u64,
            _fh: u64,
            _offset: i64,
            _size: u32,
            reply: f::ReplyData) {
        log::info!("@@@ read");
        reply.error(ENOSYS);
    }

    /// Write data.
    /// Write should return exactly the number of bytes requested except on error. An
    /// exception to this is when the file has been opened in 'direct_io' mode, in
    /// which case the return value of the write system call will reflect the return
    /// value of this operation. fh will contain the value set by the open method, or
    /// will be undefined if the open method didn't set any value.
    fn write(&mut self,
             _req: &f::Request,
             _ino: u64,
             _fh: u64,
             _offset: i64,
             _data: &[u8],
             _flags: u32,
             reply: f::ReplyWrite) {
        log::info!("@@@ write");
        reply.error(ENOSYS);
    }

    /// Flush method.
    /// This is called on each close() of the opened file. Since file descriptors can
    /// be duplicated (dup, dup2, fork), for one open call there may be many flush
    /// calls. Filesystems shouldn't assume that flush will always be called after some
    /// writes, or that if will be called at all. fh will contain the value set by the
    /// open method, or will be undefined if the open method didn't set any value.
    /// NOTE: the name of the method is misleading, since (unlike fsync) the filesystem
    /// is not forced to flush pending writes. One reason to flush data, is if the
    /// filesystem wants to return write errors. If the filesystem supports file locking
    /// operations (setlk, getlk) it should remove all locks belonging to 'lock_owner'.
    fn flush(&mut self,
             _req: &f::Request,
             _ino: u64,
             _fh: u64,
             _lock_owner: u64,
             reply: f::ReplyEmpty) {
        log::info!("@@@ flush");
        reply.error(ENOSYS);
    }

    /// Release an open file.
    /// Release is called when there are no more references to an open file: all file
    /// descriptors are closed and all memory mappings are unmapped. For every open
    /// call there will be exactly one release call. The filesystem may reply with an
    /// error, but error values are not returned to close() or munmap() which triggered
    /// the release. fh will contain the value set by the open method, or will be undefined
    /// if the open method didn't set any value. flags will contain the same flags as for
    /// open.
    fn release(&mut self,
               _req: &f::Request,
               _ino: u64,
               _fh: u64,
               _flags: u32,
               _lock_owner: u64,
               _flush: bool,
               reply: f::ReplyEmpty) {
        log::info!("@@@ release");
        reply.ok();
    }

    /// Synchronize file contents.
    /// If the datasync parameter is non-zero, then only the user data should be flushed,
    /// not the meta data.
    fn fsync(&mut self,
             _req: &f::Request,
             _ino: u64,
             _fh: u64,
             _datasync: bool,
             reply: f::ReplyEmpty) {
        log::info!("@@@ fsync");
        reply.error(ENOSYS);
    }

    /// Open a directory.
    /// Filesystem may store an arbitrary file handle (pointer, index, etc) in fh, and
    /// use this in other all other directory stream operations (readdir, releasedir,
    /// fsyncdir). Filesystem may also implement stateless directory I/O and not store
    /// anything in fh, though that makes it impossible to implement standard conforming
    /// directory stream operations in case the contents of the directory can change
    /// between opendir and releasedir.
    fn opendir(&mut self,
               _req: &f::Request,
               _ino: u64,
               _flags: u32,
               reply: f::ReplyOpen) {
        log::info!("@@@ opendir");
        log::info!("req={:?}", _req);
        log::info!("ino={:?}", _ino);
        log::info!("flags={:?}", _flags);
        reply.opened(0, 0);
    }

    /// Read directory.
    /// Send a buffer filled using buffer.fill(), with size not exceeding the
    /// requested size. Send an empty buffer on end of stream. fh will contain the
    /// value set by the opendir method, or will be undefined if the opendir method
    /// didn't set any value.
    fn readdir(&mut self,
               _req: &f::Request,
               _ino: u64,
               _fh: u64,
               _offset: i64,
               reply: f::ReplyDirectory) {
        log::info!("@@@ readdir");
        reply.error(ENOSYS);
    }

    /// Release an open directory.
    /// For every opendir call there will be exactly one releasedir call. fh will
    /// contain the value set by the opendir method, or will be undefined if the
    /// opendir method didn't set any value.
    fn releasedir(&mut self,
                  _req: &f::Request,
                  _ino: u64,
                  _fh: u64,
                  _flags: u32,
                  reply: f::ReplyEmpty) {
        log::info!("@@@ releasedir");
        reply.ok();
    }

    /// Synchronize directory contents.
    /// If the datasync parameter is set, then only the directory contents should
    /// be flushed, not the meta data. fh will contain the value set by the opendir
    /// method, or will be undefined if the opendir method didn't set any value.
    fn fsyncdir(&mut self,
                _req: &f::Request,
                _ino: u64,
                _fh: u64,
                _datasync: bool,
                reply: f::ReplyEmpty) {
        log::info!("@@@ fsyncdir");
        reply.error(ENOSYS);
    }

    /// Get file system statistics.
    fn statfs(&mut self,
              _req: &f::Request,
              _ino: u64,
              reply: f::ReplyStatfs) {
        log::info!("@@@ statfs");
        reply.statfs(0, 0, 0, 0, 0, 512, 255, 0);
    }

    /// Set an extended attribute.
    fn setxattr(&mut self,
                _req: &f::Request,
                _ino: u64,
                _name: &OsStr,
                _value: &[u8],
                _flags: u32,
                _position: u32,
                reply: f::ReplyEmpty) {
        log::info!("@@@ setxattr");
        reply.error(ENOSYS);
    }

    /// Get an extended attribute.
    /// If `size` is 0, the size of the value should be sent with `reply.size()`.
    /// If `size` is not 0, and the value fits, send it with `reply.data()`, or
    /// `reply.error(ERANGE)` if it doesn't.
    fn getxattr(&mut self,
                _req: &f::Request,
                _ino: u64,
                _name: &OsStr,
                _size: u32,
                reply: f::ReplyXattr) {
        log::info!("@@@ getxattr");
        reply.error(ENOSYS);
    }

    /// List extended attribute names.
    /// If `size` is 0, the size of the value should be sent with `reply.size()`.
    /// If `size` is not 0, and the value fits, send it with `reply.data()`, or
    /// `reply.error(ERANGE)` if it doesn't.
    fn listxattr(&mut self,
                 _req: &f::Request,
                 _ino: u64,
                 _size: u32,
                 reply: f::ReplyXattr) {
        log::info!("@@@ listtxattr");
        reply.error(ENOSYS);
    }

    /// Remove an extended attribute.
    fn removexattr(&mut self,
                   _req: &f::Request,
                   _ino: u64,
                   _name: &OsStr,
                   reply: f::ReplyEmpty) {
        log::info!("@@@ removexattr");
        reply.error(ENOSYS);
    }

    /// Check file access permissions.
    /// This will be called for the access() system call. If the 'default_permissions'
    /// mount option is given, this method is not called. This method is not called
    /// under Linux kernel versions 2.4.x
    fn access(&mut self,
              _req: &f::Request,
              _ino: u64,
              _mask: u32,
              reply: f::ReplyEmpty) {
        log::info!("@@@ access");
        reply.error(ENOSYS);
    }

    /// Create and open a file.
    /// If the file does not exist, first create it with the specified mode, and then
    /// open it. Open flags (with the exception of O_NOCTTY) are available in flags.
    /// Filesystem may store an arbitrary file handle (pointer, index, etc) in fh,
    /// and use this in other all other file operations (read, write, flush, release,
    /// fsync). There are also some flags (direct_io, keep_cache) which the
    /// filesystem may set, to change the way the file is opened. See fuse_file_info
    /// structure in <fuse_common.h> for more details. If this method is not
    /// implemented or under Linux kernel versions earlier than 2.6.15, the mknod()
    /// and open() methods will be called instead.
    fn create(&mut self,
              _req: &f::Request,
              _parent: u64,
              _name: &OsStr,
              _mode: u32,
              _flags: u32,
              reply: f::ReplyCreate) {
        log::info!("@@@ create");
        reply.error(ENOSYS);
    }

    /// Test for a POSIX file lock.
    fn getlk(&mut self,
             _req: &f::Request,
             _ino: u64,
             _fh: u64,
             _lock_owner: u64,
             _start: u64,
             _end: u64,
             _typ: u32,
             _pid: u32,
             reply: f::ReplyLock) {
        log::info!("@@@ getlk");
        reply.error(ENOSYS);
    }

    /// Acquire, modify or release a POSIX file lock.
    /// For POSIX threads (NPTL) there's a 1-1 relation between pid and owner, but
    /// otherwise this is not always the case.  For checking lock ownership,
    /// 'fi->owner' must be used. The l_pid field in 'struct flock' should only be
    /// used to fill in this field in getlk(). Note: if the locking methods are not
    /// implemented, the kernel will still allow file locking to work locally.
    /// Hence these are only interesting for network filesystems and similar.
    fn setlk(&mut self,
             _req: &f::Request,
             _ino: u64,
             _fh: u64,
             _lock_owner: u64,
             _start: u64,
             _end: u64,
             _typ: u32,
             _pid: u32,
             _sleep: bool,
             reply: f::ReplyEmpty) {
        log::info!("@@@ setlk");
        reply.error(ENOSYS);
    }

    /// Map block index within file to block index within device.
    /// Note: This makes sense only for block device backed filesystems mounted
    /// with the 'blkdev' option
    fn bmap(&mut self,
            _req: &f::Request,
            _ino: u64,
            _blocksize: u32,
            _idx: u64,
            reply: f::ReplyBmap) {
        log::info!("@@@ bmap");
        reply.error(ENOSYS);
    }
}

impl FuseFs {
    /// Mount filesystem through FUSE
    pub fn mount(mountpoint: String, client: Cl) -> Result<Session, Error> {
        let options = ["-o", "fsname=learnfs"]
            .iter()
            .map(|o| o.as_ref())
            .collect::<Vec<&OsStr>>();

        unsafe {
            let h = fuse::spawn_mount(FuseFs { client },
                                      &mountpoint, &options)?;

            Ok(h)
        }
    }
}
