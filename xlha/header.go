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
	ErrInvalidHeader = errors.New("invalid LHA header")
)

type Header struct {
	Method         string
	CompressedSize uint32
	OriginalSize   uint32
	LastModified   time.Time
	Level          uint8
	Name           string
	CRC16          uint16
}

func ReadHeader(r io.Reader) (*Header, error) {
	var headerSize uint8
	if err := binary.Read(r, binary.LittleEndian, &headerSize); err != nil {
		if err == io.EOF {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("failed to read header size: %w", err)
	}

	if headerSize == 0 {
		return nil, io.EOF
	}

	buf := make([]byte, headerSize-1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("failed to read header body: %w", err)
	}

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

	binary.Read(br, binary.LittleEndian, &h.CompressedSize)
	binary.Read(br, binary.LittleEndian, &h.OriginalSize)
	binary.Read(br, binary.LittleEndian, &dosTime)
	binary.Read(br, binary.LittleEndian, &attribute)
	binary.Read(br, binary.LittleEndian, &h.Level)

	h.LastModified = parseDosTime(dosTime)

	var nameLen uint8
	binary.Read(br, binary.LittleEndian, &nameLen)

	nameBytes := make([]byte, nameLen)
	io.ReadFull(br, nameBytes)
	h.Name = string(nameBytes)

	binary.Read(r, binary.LittleEndian, &h.CRC16)

	// Level 1 and 2 headers append extended metadata blocks (like directory paths).
	// We MUST skip these to perfectly align the stream with the compressed data.
	if h.Level == 1 {
		var osID uint8
		binary.Read(r, binary.LittleEndian, &osID)

		var extSize uint16

		// The size of the first extended header follows the OS-ID, on the raw stream.
		if err := binary.Read(r, binary.LittleEndian, &extSize); err == nil {

			// Chain-read and discard extended headers until we hit a 0x0000 terminator
			for extSize != 0 {
				if extSize >= 2 {
					// Discard the extended header data (extSize includes the 2-byte size itself)
					if _, err := io.CopyN(io.Discard, r, int64(extSize-2)); err != nil {
						break
					}
				}

				// Read the size of the NEXT extended header directly from the raw stream
				if err := binary.Read(r, binary.LittleEndian, &extSize); err != nil {
					break
				}
			}
		}
	}

	return h, nil
}

func parseDosTime(t uint32) time.Time {
	date := t >> 16
	timePart := t & 0xFFFF

	year := int((date >> 9) + 1980)
	month := time.Month((date >> 5) & 0x0F)
	day := int(date & 0x1F)

	hour := int(timePart >> 11)
	minute := int((timePart >> 5) & 0x3F)
	second := int((timePart & 0x1F) * 2)

	return time.Date(year, month, day, hour, minute, second, 0, time.UTC)
}
