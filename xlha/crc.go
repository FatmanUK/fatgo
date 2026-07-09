package xlha

import (
	"hash"
)

const (
	// The standard polynomial used by LHA, ARC, and LZH formats.
	// It is the bit-reversed version of the standard CRC-16-IBM polynomial (0x8005).
	crcPolynomial = 0xA001
)

// crcTable caches the precomputed values for byte-by-byte CRC calculation.
var crcTable [256]uint16

// init automatically calculates the lookup table the first time the package is imported.
// This gives us the performance of a hardcoded C array without cluttering the code.
func init() {
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ crcPolynomial
			} else {
				crc >>= 1
			}
		}
		crcTable[i] = crc
	}
}

// crc16 implements the standard library's hash.Hash16 interface.
type crc16 struct {
	val uint16
}

// NewCRC16 returns a new hash.Hash16 computing the LHA-specific CRC-16 checksum.
func NewCRC16() hash.Hash16 {
	return &crc16{val: 0}
}

// Write processes a slice of bytes and updates the running checksum.
// This allows the CRC calculator to be used seamlessly with io.Writers.
func (c *crc16) Write(p []byte) (n int, err error) {
	for _, b := range p {
		// Calculate the next state using the precomputed lookup table
		c.val = crcTable[byte(c.val)^b] ^ (c.val >> 8)
	}
	return len(p), nil
}

// Sum16 returns the current 16-bit checksum value.
func (c *crc16) Sum16() uint16 {
	return c.val
}

// Sum appends the current checksum to the provided byte slice.
func (c *crc16) Sum(b []byte) []byte {
	s := c.Sum16()
	return append(b, byte(s>>8), byte(s))
}

// Reset clears the checksum back to its initial state (0x0000 for LHA).
func (c *crc16) Reset() {
	c.val = 0
}

// Size returns the number of bytes Sum will return (2 bytes for CRC-16).
func (c *crc16) Size() int {
	return 2
}

// BlockSize returns the hash's underlying block size.
func (c *crc16) BlockSize() int {
	return 1
}
