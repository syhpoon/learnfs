pub struct Inode {
    // File type
    ftype: u16,
    // File mode
    mode: u16,
    // Number of directories containing this inode
    nlinks: u16,
    // Owner's user id
    uid: u16,
    // Owner's group id
    gid: u16,
    // File size in bytes
    size: u32,
    // Access time
    atime: u32,
    // Modification time
    mtime: u32,
    // Status change time
    ctime: u32,
    // Direct block pointers
    blocks: [u32; 7],
    indirect_block: u32,
    double_indirect_block: u32,
    padding: [u8; 4],
}
