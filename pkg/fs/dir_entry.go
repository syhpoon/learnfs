// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"fmt"
	"io"
)

const (
	MAX_FILE_NAME  = 60
	DIR_ENTRY_SIZE = 64
)

type DirEntry struct {
	Name     [MAX_FILE_NAME]byte
	InodePtr InodePtr

	nameLen int32
}

func NewDirEntry(inodePtr InodePtr, name string) (*DirEntry, error) {
	l := len(name)
	if l > MAX_FILE_NAME {
		return nil, fmt.Errorf("file name too long")
	}

	de := &DirEntry{
		InodePtr: inodePtr,
		Name:     [60]byte{},
		nameLen:  int32(l),
	}

	copy(de.Name[:], name)

	return de, nil
}

func (de *DirEntry) GetName() string {
	return string(de.Name[:de.nameLen])
}

func (de *DirEntry) EncodeTo(w io.Writer) error {
	fields := []any{
		de.Name,
		de.InodePtr,
	}

	return encodeFields(w, fields)
}

func (de *DirEntry) DecodeFrom(r io.Reader) error {
	fields := []any{
		&de.Name,
		&de.InodePtr,
	}

	if err := decodeFields(r, fields); err != nil {
		return err
	}

	nameLen := int32(0)
	for i := range de.Name {
		if de.Name[i] == 0 {
			break
		}

		nameLen++
	}

	de.nameLen = nameLen

	return nil
}
