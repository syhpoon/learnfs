// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"encoding/binary"
	"io"
)

var byteOrder = binary.LittleEndian

type Encoder interface {
	EncodeTo(io.Writer) error
}

type Decoder interface {
	DecodeFrom(io.Reader) error
}

func encodeFields(w io.Writer, fields []any) error {
	for i := range fields {
		if err := binary.Write(w, byteOrder, fields[i]); err != nil {
			return err
		}
	}

	return nil
}

func decodeFields(r io.Reader, fields []any) error {
	for i := range fields {
		if err := binary.Read(r, byteOrder, fields[i]); err != nil {
			return err
		}
	}

	return nil
}
