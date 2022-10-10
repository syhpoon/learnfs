// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"syscall"
)

type ErrorSystem struct {
	Errno syscall.Errno
}

func (err *ErrorSystem) Error() string {
	return err.Errno.Error()
}

var (
	ErrorNotFound = &ErrorSystem{Errno: syscall.ENOENT}
	ErrorExists   = &ErrorSystem{Errno: syscall.EEXIST}
)
