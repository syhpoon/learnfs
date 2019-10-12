/*
 MIT License

 Copyright (c) 2019 Max Kuznetsov <syhpoon@syhpoon.ca>

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 SOFTWARE.
*/

#[derive(PartialEq, Debug, Default)]
pub struct Bitmap {
    buf: Vec<u8>,
    max_bit: usize,
}

impl Bitmap {
    pub fn new(size: usize) -> Self {
        assert!(size > 0, "bitmap size cannot be zero");

        Bitmap {
            buf: vec![0; size],
            max_bit: size * 8,
        }
    }

    pub fn from_buf(buf: Vec<u8>) -> Self {
        let size = buf.len();

        assert!(size > 0, "bitmap size cannot be zero");

        Bitmap {
            buf,
            max_bit: size * 8,
        }
    }

    pub fn as_slice(&self) -> &[u8] {
        return self.buf.as_slice();
    }

    pub fn set(&mut self, bit: usize) {
        let (idx, off) = self.convert(bit);

        self.buf[idx] |= 1 << off
    }

    pub fn clear(&mut self, bit: usize) {
        let (idx, off) = self.convert(bit);

        self.buf[idx] &= !(1 << off)
    }

    pub fn is_set(&self, bit: usize) -> bool {
        let (idx, off) = self.convert(bit);

        return (self.buf[idx] >> off) & 1 > 0;
    }

    /// Find the next free bit number searching from the given index
    pub fn next_clear_bit(&mut self, from: usize) -> Option<usize> {
        let mut scanned = 0;
        let mut idx = from;

        while scanned <= self.max_bit {
            scanned += 1;

            if !self.is_set(idx) {
                return Some(idx);
            }

            idx = (idx + 1) % self.max_bit;
        }

        None
    }

    /// Convert raw bit number into vector index and an offset
    #[inline(always)]
    fn convert(&self, bit: usize) -> (usize, u8) {
        let idx = bit / 8;
        let off = (bit - (idx * 8)) as u8;

        (idx, off)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_bitmap1() {
        let mut b = Bitmap::new(1);
        let cases = vec![
            (0, 1u8),
            (1, 2u8),
            (2, 4u8),
            (3, 8u8),
            (4, 16u8),
            (5, 32u8),
            (6, 64u8),
            (7, 128u8),
        ];

        assert_eq!(b.buf[0], 0u8);

        for (bit, val) in cases {
            b.set(bit);
            assert_eq!(b.buf[0], val);
            b.clear(bit);
        }

        assert_eq!(b.buf[0], 0u8);
    }

    #[test]
    fn test_bitmap2() {
        let mut b = Bitmap::new(1);

        assert_eq!(b.buf[0], 0u8);

        for i in 0..8 {
            b.set(i);
        }

        assert_eq!(b.buf[0], 255u8);

        for i in 0..8 {
            b.clear(i);
        }

        assert_eq!(b.buf[0], 0u8);
    }

    #[test]
    fn test_bitmap3() {
        let mut b = Bitmap::new(4);

        b.set(26);
        assert!(b.is_set(26));

        assert_eq!(b.buf[0], 0u8);
        assert_eq!(b.buf[1], 0u8);
        assert_eq!(b.buf[2], 0u8);
        assert_eq!(b.buf[3], 4u8);
    }

    #[test]
    fn test_bitmap_find_next() {
        let mut b = Bitmap::new(1);

        assert_eq!(b.next_clear_bit(0), Some(0));

        b.set(0);
        b.set(1);
        b.set(2);
        b.set(3);
        b.set(4);

        assert_eq!(b.next_clear_bit(0), Some(5));

        b.set(6);
        b.set(7);

        assert_eq!(b.next_clear_bit(6), Some(5));

        b.set(5);
        assert_eq!(b.next_clear_bit(0), None);
    }
}
