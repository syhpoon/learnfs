// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package device

import (
	"context"
	"fmt"
	"sync"
)

// In-memory device, used for testing
type Memory struct {
	size uint
	data []byte
	ctx  context.Context

	sync.RWMutex
}

var _ Device = &Memory{}

func NewMemory(ctx context.Context, size uint) *Memory {
	dev := &Memory{
		size: size,
		data: make([]byte, size),
		ctx:  ctx,
	}

	return dev
}

func (mem *Memory) Capacity() uint64 {
	return uint64(mem.size)
}

func (mem *Memory) Read(ops ...*Op) error {
	mem.RLock()
	defer mem.RUnlock()

	for _, op := range ops {
		if err := mem.checkSize(op); err != nil {
			return err
		}

		copy(op.Buf, mem.data[op.Offset:int(op.Offset)+len(op.Buf)])
	}

	return nil
}

func (mem *Memory) Write(ops ...*Op) error {
	mem.Lock()
	defer mem.Unlock()

	for _, op := range ops {
		if err := mem.checkSize(op); err != nil {
			return err
		}

		copy(mem.data[op.Offset:int(op.Offset)+len(op.Buf)], op.Buf)
	}

	return nil
}

func (mem *Memory) Close() error {
	return nil
}

func (mem *Memory) checkSize(op *Op) error {
	s := uint(op.Offset) + uint(len(op.Buf))

	if s > mem.size {
		return fmt.Errorf("offset+buffer size (%d) exceeding device capacity %d", s, mem.size)
	}

	return nil
}
