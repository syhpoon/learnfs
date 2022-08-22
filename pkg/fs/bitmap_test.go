// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import (
	"testing"

	"github.com/barweiss/go-tuple"
	"github.com/stretchr/testify/require"
)

func TestBitmap1(t *testing.T) {
	b := NewBitmap(7, make([]byte, 1))

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

		b.Set(bit)
		require.Equal(t, val, b.buf[0])
		b.Clear(bit)
	}

	require.Equal(t, uint8(0), b.buf[0])
}

func TestBitmap2(t *testing.T) {
	b := NewBitmap(7, make([]byte, 1))

	require.Equal(t, uint8(0), b.buf[0])

	for i := uint32(0); i < 8; i++ {
		b.Set(i)
	}

	require.Equal(t, uint8(255), b.buf[0])

	for i := uint32(0); i < 8; i++ {
		b.Clear(i)
	}

	require.Equal(t, uint8(0), b.buf[0])
}

func TestBitmap3(t *testing.T) {
	b := NewBitmap(512, make([]byte, 4))

	b.Set(26)
	require.True(t, b.IsSet(26))

	require.Equal(t, uint8(0), b.buf[0])
	require.Equal(t, uint8(0), b.buf[1])
	require.Equal(t, uint8(0), b.buf[2])
	require.Equal(t, uint8(4), b.buf[3])
}

func TestBitmapFindNext(t *testing.T) {
	b := NewBitmap(7, make([]byte, 1))

	require.Equal(t, *b.NextClearBit(0), uint32(0))

	b.Set(0)
	b.Set(1)
	b.Set(2)
	b.Set(3)
	b.Set(4)
	require.Equal(t, *b.NextClearBit(0), uint32(5))

	b.Set(6)
	b.Set(7)
	require.Equal(t, *b.NextClearBit(0), uint32(5))

	b.Set(5)
	require.Equal(t, b.NextClearBit(0), (*uint32)(nil))
}
