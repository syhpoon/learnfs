// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package device

type Op struct {
	Buf    []byte
	Offset uint64
}

type Device interface {
	Read(...*Op) error
	Write(...*Op) error
	Capacity() uint64
	Close() error
}
