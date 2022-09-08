Overview
========

`LearnFS` is a toy, user-space network filesystem. It was made to learn basic
concepts of the filesystem design. Despite implemented entirely in user-space,
`LearnFS` nevertheless employs many concepts, algorithms and approaches used
in normal filesystems, mainly Unix/Linux-based ones, such as Minix Filesystem,
Ext2 and FFS.

A lot of efforts have been made to keep the implementation as simple as possible,
to allow people, interested in filesystem design to easily follow and understand
the code. Because of this, `LearnFS` is not suitable for any use case where
data loss is impactful and performance is important.

LearnFS in a nutshell
---------------------
* Implemented entirely is user-space in Go programming language.
* Using client-server model, with one server and many concurrent clients,
  potentially running on different machines.
* Can be created either inside a pre-allocated file in an existing filesystem,
  or directly on a block device using [io_uring](https://kernel.dk/io_uring.pdf)
  Linux interface.
* Client-server communication is done with [gRPC](https://grpc.io/).

Architecture
------------

![Architecture](docs/architecture.png)

Getting started
---------------
TODO