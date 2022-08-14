// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"bytes"
	"fmt"
	"sync"
)

type BufPool struct {
	pool      *sync.Pool
	blockSize uint32
}

func NewBufPool(blockSize uint32) *BufPool {
	pool := &sync.Pool{
		New: func() any {
			return make(Buf, blockSize)
		},
	}

	return &BufPool{
		pool:      pool,
		blockSize: blockSize,
	}
}

func (bp *BufPool) GetRead() Buf {
	return bp.pool.Get().(Buf)
}

func (bp *BufPool) GetWrite() Buf {
	return bp.pool.Get().(Buf)[:0]
}

func (bp *BufPool) Put(buf Buf) {
	bp.pool.Put(buf[:cap(buf)])
}

func serialize(pool *BufPool, objects ...any) (Buf, error) {
	buf := pool.GetWrite()
	bbuf := bytes.NewBuffer(buf)

	for i := range objects {
		switch obj := objects[i].(type) {
		case []byte:
			if _, err := bbuf.Write(obj); err != nil {
				return nil, fmt.Errorf("failed to write bytes object to buffer: %w", err)
			}
		case Encoder:
			if err := obj.EncodeTo(bbuf); err != nil {
				return nil, fmt.Errorf("failed to encode object to buffer: %w", err)
			}
		default:
			return nil, fmt.Errorf("unable to serialize type %T", obj)
		}
	}

	mod := uint32(bbuf.Cap()) % pool.blockSize
	if mod != 0 {
		bbuf.Grow(int(pool.blockSize - mod))
	}

	// Write padding if needed
	if bbuf.Len() < bbuf.Cap() {
		// TODO: make a pool of paddings?
		// TODO: or pre-initialize or possible size from 1 to block-size
		pad := make([]byte, bbuf.Cap()-(bbuf.Len()))

		if _, err := bbuf.Write(pad); err != nil {
			return nil, fmt.Errorf("failed to write padding to buffer: %w", err)
		}
	}

	return bbuf.Bytes(), nil

}
