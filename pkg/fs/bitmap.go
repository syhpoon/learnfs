// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

type bitmap struct {
	buf    []byte
	maxBit uint32
}

func newBitmap(size uint32, buf []byte) *bitmap {
	return &bitmap{
		buf:    buf,
		maxBit: size,
	}
}

func (bm *bitmap) size() int {
	return len(bm.buf)
}

func (bm *bitmap) set(bit uint32) {
	idx, off := bm.convert(bit)

	bm.buf[idx] |= 1 << off
}

func (bm *bitmap) clear(bit uint32) {
	idx, off := bm.convert(bit)

	bm.buf[idx] &= ^(1 << off)
}

func (bm *bitmap) isSet(bit uint32) bool {
	idx, off := bm.convert(bit)

	return (bm.buf[idx]>>off)&1 > 0
}

// Find the next free bit number searching from the given index
func (bm *bitmap) nextClearBit(from uint32) *uint32 {
	scanned := uint32(0)
	idx := from

	for scanned <= bm.maxBit {
		scanned += 1

		if !bm.isSet(idx) {
			return &idx
		}

		idx = (idx + 1) % bm.maxBit
	}

	return nil
}

// Convert raw bit number into vector index and an offset
func (bm *bitmap) convert(bit uint32) (uint32, uint8) {
	idx := bit / 8
	off := bit - (idx * 8)

	return idx, uint8(off)
}

// Return the number of free bits
func (bm *bitmap) numFreeBits() uint32 {
	count := uint32(0)

	for i := uint32(0); i < bm.maxBit; i++ {
		if !bm.isSet(i) {
			count++
		}
	}

	return count
}
