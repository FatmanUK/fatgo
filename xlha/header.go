package xlha

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

var (
	// ErrInvalidHeader indicates the header data is corrupt or not an LHA archive.
	ErrInvalidHeader = errors.New("invalid LHA header")
)

// Header represents the extracted metadata for a file within the archive.
type Header struct {
	Method         string    // e.g., "-lh5-", "-lh0-"
	CompressedSize uint32    // Size of the compressed data stream
	OriginalSize   uint32    // Size of the file when extracted
	LastModified   time.Time // Converted from MS-DOS timestamp
	Level          uint8     // Header level (0, 1, or 2)
	Name           string    // Filename
	CRC16          uint16    // Expected CRC16 checksum of the uncompressed file
}

// ReadHeader reads and parses an LHA header from the provided stream.
func ReadHeader(r io.Reader) (*Header, error) {
	// The first byte of an LHA header is the header size.
	// A size of 0 indicates the end of the archive.
	var headerSize uint8
	if err := binary.Read(r, binary.LittleEndian, &headerSize); err != nil {
		if err == io.EOF {
			return nil, io.EOF // End of archive gracefully reached
		}
		return nil, fmt.Errorf("failed to read header size: %w", err)
	}

	if headerSize == 0 {
		return nil, io.EOF
	}

	// Read the rest of the base header into a buffer.
	// We subtract 1 because we already read the headerSize byte.
	buf := make([]byte, headerSize-1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("failed to read header body: %w", err)
	}

	// Use bytes.Reader for easy, panic-free sequential reading from our buffer
	br := bytes.NewReader(buf)

	var checksum uint8
	binary.Read(br, binary.LittleEndian, &checksum)

	methodBytes := make([]byte, 5)
	io.ReadFull(br, methodBytes)

	h := &Header{
		Method: string(methodBytes),
	}

	var dosTime uint32
	var attribute uint8

	// Read the fixed-size numerical fields (Little Endian is standard for LHA)
	binary.Read(br, binary.LittleEndian, &h.CompressedSize)
	binary.Read(br, binary.LittleEndian, &h.OriginalSize)
	binary.Read(br, binary.LittleEndian, &dosTime)
	binary.Read(br, binary.LittleEndian, &attribute)
	binary.Read(br, binary.LittleEndian, &h.Level)

	h.LastModified = parseDosTime(dosTime)

	// Filename length and filename
	var nameLen uint8
	binary.Read(br, binary.LittleEndian, &nameLen)

	nameBytes := make([]byte, nameLen)
	io.ReadFull(br, nameBytes)
	h.Name = string(nameBytes)

	// Read CRC16
	binary.Read(br, binary.LittleEndian, &h.CRC16)

	// Note: For Level 1 and Level 2 headers, there are "Extended Headers"
	// that follow. We can add logic to skip or parse those next if needed.

	return h, nil
}

// parseDosTime converts a 32-bit MS-DOS timestamp to a Go time.Time object.
// This replaces the bitwise macros often found in the C implementation.
func parseDosTime(t uint32) time.Time {
	date := t >> 16
	timePart := t & 0xFFFF

	year := int((date >> 9) + 1980)
	month := time.Month((date >> 5) & 0x0F)
	day := int(date & 0x1F)

	hour := int(timePart >> 11)
	minute := int((timePart >> 5) & 0x3F)
	second := int((timePart & 0x1F) * 2) // DOS seconds have a 2-second resolution

	return time.Date(year, month, day, hour, minute, second, 0, time.UTC)
}
