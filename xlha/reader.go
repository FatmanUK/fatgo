package xlha

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrCRCMismatch is returned when the computed checksum does not match the header.
	ErrCRCMismatch = errors.New("lha: CRC checksum mismatch")

	// ErrUnsupportedMethod is returned when the archive uses a compression scheme other than -lh0- or -lh5-.
	ErrUnsupportedMethod = errors.New("lha: unsupported compression method")
)

// Reader provides sequential access to the contents of an LHA archive.
// It implements the io.Reader interface to stream uncompressed file contents.
type Reader struct {
	r             io.Reader
	header        *Header
	decoder       io.Reader
	limitedSource *io.LimitedReader // Tracks unconsumed compressed bytes
	crc           Hash16            // Updated to use the local Hash16 interface
	bytesLeft     int64             // Tracks remaining uncompressed bytes for the current file
}

// NewReader creates a new Reader reading from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// Next advances to the next entry in the LHA archive and returns its Header.
func (tr *Reader) Next() (*Header, error) {
	// 1. If a file was partially read (or skipped entirely), consume
	// any remaining compressed data to align the stream with the next header block.
	if tr.limitedSource != nil && tr.limitedSource.N > 0 {
		if _, err := io.CopyN(io.Discard, tr.r, tr.limitedSource.N); err != nil {
			return nil, fmt.Errorf("lha: failed to skip remaining file data: %w", err)
		}
	}

	// 2. Read the next metadata header block from the stream
	h, err := ReadHeader(tr.r)
	if err != nil {
		return nil, err
	}

	tr.header = h
	tr.bytesLeft = int64(h.OriginalSize)
	tr.crc = NewCRC16()

	// 3. Enforce strict reading limits based on the header's declared CompressedSize.
	tr.limitedSource = &io.LimitedReader{R: tr.r, N: int64(h.CompressedSize)}

	// 4. Instantiated processing pipeline based on the compression method identifier
	switch h.Method {
	case "-lh0-":
		tr.decoder = tr.limitedSource
	case "-lh5-":
		tr.decoder = NewLH5Decoder(tr.limitedSource, h.OriginalSize)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedMethod, h.Method)
	}

	return h, nil
}

// Read reads uncompressed data from the current file entry in the LHA archive.
func (tr *Reader) Read(b []byte) (int, error) {
	if tr.decoder == nil {
		return 0, io.EOF
	}

	// Edge-case: Handle zero-byte files smoothly while confirming their CRC
	if tr.bytesLeft == 0 {
		if tr.crc != nil {
			if tr.crc.Sum16() != tr.header.CRC16 {
				tr.crc = nil
				return 0, ErrCRCMismatch
			}
			tr.crc = nil
		}
		return 0, io.EOF
	}

	// Run the read operation through the selected decompression layout
	n, err := tr.decoder.Read(b)
	if n > 0 {
		// Feed the uncompressed data into our CRC-16 engine
		_, _ = tr.crc.Write(b[:n])
		tr.bytesLeft -= int64(n)

		// Verification trigger: We have collected exactly the number of uncompressed bytes expected
		if tr.bytesLeft == 0 {
			if tr.crc.Sum16() != tr.header.CRC16 {
				return n, ErrCRCMismatch
			}
			return n, io.EOF
		}
	}

	// Catch scenario where decompression ends prematurely before extracting all original bytes
	if err == io.EOF && tr.bytesLeft > 0 {
		return n, io.ErrUnexpectedEOF
	}

	return n, err
}
