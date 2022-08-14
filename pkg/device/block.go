// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package device

/*
	#include <linux/fs.h>
*/
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const BLKGETSIZE64 = C.BLKGETSIZE64

func blkGetSize64(path string) (uint64, error) {
	var size uint64

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), BLKGETSIZE64, uintptr(unsafe.Pointer(&size))); err != 0 {
		return 0, fmt.Errorf("failed to call ioctl BLKGETSIZE64: %w", err)
	}

	return size, nil
}
