// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"testing"

	"github.com/barweiss/go-tuple"
	"github.com/stretchr/testify/require"
)

func TestBitmap1(t *testing.T) {
	b := newBitmap(7, make([]byte, 1))

	cases := []tuple.T2[uint32, uint8]{
		tuple.New2[uint32, uint8](0, 1),
		tuple.New2[uint32, uint8](1, 2),
		tuple.New2[uint32, uint8](2, 4),
		tuple.New2[uint32, uint8](3, 8),
		tuple.New2[uint32, uint8](4, 16),
		tuple.New2[uint32, uint8](5, 32),
		tuple.New2[uint32, uint8](6, 64),
		tuple.New2[uint32, uint8](7, 128),
	}

	require.Equal(t, uint8(0), b.buf[0])

	for _, tup := range cases {
		bit := tup.V1
		val := tup.V2

		b.set(bit)
		require.Equal(t, val, b.buf[0])
		b.clear(bit)
	}

	require.Equal(t, uint8(0), b.buf[0])
}

func TestBitmap2(t *testing.T) {
	b := newBitmap(7, make([]byte, 1))

	require.Equal(t, uint8(0), b.buf[0])

	for i := uint32(0); i < 8; i++ {
		b.set(i)
	}

	require.Equal(t, uint8(255), b.buf[0])

	for i := uint32(0); i < 8; i++ {
		b.clear(i)
	}

	require.Equal(t, uint8(0), b.buf[0])
}

func TestBitmap3(t *testing.T) {
	b := newBitmap(512, make([]byte, 4))

	b.set(26)
	require.True(t, b.isSet(26))

	require.Equal(t, uint8(0), b.buf[0])
	require.Equal(t, uint8(0), b.buf[1])
	require.Equal(t, uint8(0), b.buf[2])
	require.Equal(t, uint8(4), b.buf[3])
}

func TestBitmapFindNext(t *testing.T) {
	b := newBitmap(7, make([]byte, 1))

	require.Equal(t, *b.nextClearBit(0), uint32(0))

	b.set(0)
	b.set(1)
	b.set(2)
	b.set(3)
	b.set(4)
	require.Equal(t, *b.nextClearBit(0), uint32(5))

	b.set(6)
	b.set(7)
	require.Equal(t, *b.nextClearBit(0), uint32(5))

	b.set(5)
	require.Equal(t, b.nextClearBit(0), (*uint32)(nil))
}
