// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

type Bitmap struct {
	buf    []byte
	maxBit uint32
}

func NewBitmap(size uint32, buf []byte) *Bitmap {
	return &Bitmap{
		buf:    buf,
		maxBit: size,
	}
}

func (bm *Bitmap) Size() int {
	return len(bm.buf)
}

func (bm *Bitmap) Set(bit uint32) {
	idx, off := bm.convert(bit)

	bm.buf[idx] |= 1 << off
}

func (bm *Bitmap) Clear(bit uint32) {
	idx, off := bm.convert(bit)

	bm.buf[idx] &= ^(1 << off)
}

func (bm *Bitmap) IsSet(bit uint32) bool {
	idx, off := bm.convert(bit)

	return (bm.buf[idx]>>off)&1 > 0
}

func (bm *Bitmap) GetBuf() []byte {
	return bm.buf
}

// Find the next free bit number searching from the given index
func (bm *Bitmap) NextClearBit(from uint32) *uint32 {
	scanned := uint32(0)
	idx := from

	for scanned <= bm.maxBit {
		scanned += 1

		if !bm.IsSet(idx) {
			return &idx
		}

		idx = (idx + 1) % bm.maxBit
	}

	return nil
}

// Convert raw bit number into vector index and an offset
func (bm *Bitmap) convert(bit uint32) (uint32, uint8) {
	idx := bit / 8
	off := bit - (idx * 8)

	return idx, uint8(off)
}
