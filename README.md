# Overview

`LearnFS` is a toy, user-space network filesystem. It was made to learn basic
concepts of the filesystem design. Despite implemented entirely in user-space,
`LearnFS` nevertheless employs many concepts, algorithms and approaches used
in normal filesystems, mainly Unix/Linux-based ones, such as Minix Filesystem,
Ext2 and FFS.

A lot of efforts have been made to keep the implementation as simple as possible,
to allow people, interested in filesystem design to easily follow and understand
the code. Because of this, `LearnFS` is not suitable for any use case where
data loss is impactful and performance is important.

# LearnFS in a nutshell

* Implemented entirely is user-space in Go programming language.
* Using client-server model, with one server and many concurrent clients,
  potentially running on different machines.
* Can be created either inside a pre-allocated file in an existing filesystem,
  or directly on a block device using [io_uring](https://kernel.dk/io_uring.pdf)
  Linux interface.
* Client-server communication is done with [gRPC](https://grpc.io/).
* 4GB max file size

# Architecture

![Architecture](docs/architecture.png)

# Installation

```shell
go install github.com/syhpoon/learnfs@latest
```

# Getting started

## Select a device

First, decide on a device to use for the `LearnFS`. It currently supports both 
regular files on existing filesystem or raw block device. Using raw block
device requires root privilege, though.

The easiest way is to allocate an empty file:

```shell
$ fallocate -l 100m ~/learnfs
```

## Create a filesystem

```shell
$ learnfs mkfs ~/learnfs
```

## Check the fs info
```shell
$ learnfs info ~/learnfs 
{
  "Magic": 32482246873146732,
  "BlockSize": 4096,
  "NumInodes": 6400,
  "InodeSize": 128,
  "NumInodeBitmapBlocks": 1,
  "NumDataBitmapBlocks": 1,
  "NumInodeBlocks": 200,
  "NumDataBlocks": 25398,
  "FirstInodeBlockOffset": 9216,
  "FirstDataBlockOffset": 828416
}
```

## Run LearnFS server

```shell
$ learnfs server ~/learnfs 
2022-10-23T22:51:01-04:00 INF pkg/server/server.go:42 > starting learnfs server listen=:1234
```

## Run LearnFS client

In a separate shell:
```shell
mkdir /tmp/learnfs/ && learnfs mount /tmp/learnfs/
2022-10-23T22:53:29-04:00 DBG pkg/fuse/req_getattr.go:17 > getAttr req={"Flags":0,"Handle":0}
2022-10-23T22:53:29-04:00 DBG pkg/fuse/req_lookup.go:16 > lookup req={"Name":".Trash"}
2022-10-23T22:53:29-04:00 DBG pkg/fuse/req_lookup.go:16 > lookup req={"Name":".Trash-1000"}
```

## Play with mounted LearnFS at /tmp/learnfs

```shell
$ df -h /tmp/learnfs/
Filesystem      Size  Used Avail Use% Mounted on
learnfs         100M  8.0K  100M   1% /tmp/learnfs
$ echo LearnFS > /tmp/learnfs/file
$ cat /tmp/learnfs/file
LearnFS
```