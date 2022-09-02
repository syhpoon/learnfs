// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package device

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/dshulyak/uring"
	"github.com/samber/mo"
)

type BlockUring struct {
	ring     *uring.Ring
	f        *os.File
	fd       uintptr
	capacity uint64
	id       uint64
	res      map[uint64]chan mo.Result[int32]
	ctx      context.Context
	cancel   context.CancelFunc

	sync.Mutex
}

var _ Device = &BlockUring{}

func NewBlockUring(ctx context.Context, path string) (*BlockUring, error) {
	// Get device capacity
	size, err := getDeviceSize(path)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_DIRECT, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open device %s: %w", path, err)
	}

	// TODO: Consider io polling mode at some point?
	ring, err := uring.Setup(16, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create io_uring instance: %w", err)
	}

	if err := ring.RegisterFiles([]int32{int32(f.Fd())}); err != nil {
		return nil, fmt.Errorf("failed to register file descriptor: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	dev := &BlockUring{
		ring:     ring,
		f:        f,
		fd:       f.Fd(),
		capacity: size,
		id:       0,
		res:      map[uint64]chan mo.Result[int32]{},
		ctx:      ctx,
		cancel:   cancel,
	}

	go dev.poll()

	return dev, nil
}

func (bu *BlockUring) Capacity() uint64 {
	return bu.capacity
}

func (bu *BlockUring) Read(ops ...*Op) error {
	return bu.run(ops, func(sqe *uring.SQEntry, opIdx int) {
		uring.Read(sqe, 0, ops[opIdx].Buf)
	})
}

func (bu *BlockUring) Write(ops ...*Op) error {
	return bu.run(ops, func(sqe *uring.SQEntry, opIdx int) {
		uring.Write(sqe, 0, ops[opIdx].Buf)
	})
}

func (bu *BlockUring) getId() uint64 {
	return atomic.AddUint64(&bu.id, 1)
}

func (bu *BlockUring) isAligned(val uint64) bool {
	return val%uint64(512) == 0
}

func (bu *BlockUring) run(ops []*Op, cb func(*uring.SQEntry, int)) error {
	id := bu.getId()
	ch := make(chan mo.Result[int32], len(ops))

	bu.Lock()
	bu.res[id] = ch
	bu.Unlock()

	for i := range ops {
		if !bu.isAligned(ops[i].Offset) {
			return fmt.Errorf("offset is not a multiple of 512: %d", ops[i].Offset)
		}

		if !bu.isAligned(uint64(len(ops[i].Buf))) {
			return fmt.Errorf("buf size is not a multiple of 512: %d", len(ops[i].Buf))
		}

		sqe := bu.ring.GetSQEntry()

		sqe.SetFlags(uring.IOSQE_FIXED_FILE)
		sqe.SetUserData(id)
		sqe.SetOffset(ops[i].Offset)

		cb(sqe, i)

		if _, err := bu.ring.Submit(0); err != nil {
			return fmt.Errorf("failed to submit ops: %w", err)
		}
	}

	var err error

	// Wait for responses
	for i := len(ops); i > 0; i-- {
		res := <-ch

		if res.IsError() {
			err = fmt.Errorf("failed to perform op: %s", res.Error())
		} else {
			code, _ := res.Get()
			if code < 0 {
				errno := syscall.Errno(-code)
				err = fmt.Errorf("op returned error: %s (%d)", errno.Error(), code)
			}
		}
	}

	bu.Lock()
	delete(bu.res, id)
	bu.Unlock()

	return err
}

func (bu *BlockUring) Close() error {
	bu.cancel()

	_ = bu.ring.Close()
	_ = bu.f.Close()

	return nil
}

func (bu *BlockUring) poll() {
	for {
		cqe, err := bu.ring.GetCQEntry(1)
		if err == syscall.EINTR {
			continue
		}

		var r mo.Result[int32]
		if err != nil {
			r = mo.Err[int32](err)
		} else {
			r = mo.Ok[int32](cqe.Result())
		}

		bu.Lock()
		ch := bu.res[cqe.UserData()]
		bu.Unlock()

		ch <- r
	}
}

func getDeviceSize(path string) (uint64, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(path, &stat)
	if err != nil {
		return 0, fmt.Errorf("failed to stat device: %w", err)
	}

	switch stat.Mode & syscall.S_IFMT {
	case syscall.S_IFREG:
		return uint64(stat.Size), nil
	case syscall.S_IFBLK:
		return blkGetSize64(path)
	default:
		return 0, fmt.Errorf("only regular files and block devices are supported")
	}
}
