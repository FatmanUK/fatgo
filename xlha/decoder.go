package xlha

import (
	"bufio"
	"errors"
	"io"
)

const (
	// LH5 uses an 8KB sliding window for its dictionary
	windowSize = 8192
)

// ==========================================
// PART 1: THE BIT READER
// ==========================================

// BitReader allows reading an arbitrary number of bits from a byte stream.
type BitReader struct {
	r    io.ByteReader
	bits uint16 // Buffer holding up to 16 bits
	n    uint8  // Number of valid bits currently in the buffer
}

// ReadBits reads up to 16 bits from the underlying stream.
func (b *BitReader) ReadBits(count uint8) (uint16, error) {
	for b.n < count {
		nextByte, err := b.r.ReadByte()
		if err != nil {
			return 0, err
		}
		// LHA packs bits starting from the most significant bit
		b.bits = (b.bits << 8) | uint16(nextByte)
		b.n += 8
	}

	// Extract the requested bits
	res := b.bits >> (b.n - count)
	
	// Create a mask to clear the bits we just read
	mask := (uint16(1) << (b.n - count)) - 1
	b.bits &= mask
	b.n -= count

	return res, nil
}

// ==========================================
// PART 2: THE HUFFMAN TREE
// ==========================================

// huffmanTree represents a canonical Huffman decoding tree.
type huffmanTree struct {
	left     []uint16
	right    []uint16
	maxNodes uint16
}

// newHuffmanTree initializes the slices needed for the tree.
func newHuffmanTree(numSymbols uint16) *huffmanTree {
	maxNodes := numSymbols * 2
	return &huffmanTree{
		left:     make([]uint16, maxNodes),
		right:    make([]uint16, maxNodes),
		maxNodes: maxNodes,
	}
}

// buildTree reconstructs the Huffman tree from an array of code lengths.
func (t *huffmanTree) buildTree(lengths []uint8) error {
	numSymbols := uint16(len(lengths))
	
	count := make([]uint16, 17)
	for _, length := range lengths {
		if length > 16 {
			return errors.New("invalid Huffman code length")
		}
		count[length]++
	}

	startCode := make([]uint16, 17)
	code := uint16(0)
	for i := uint8(1); i <= 16; i++ {
		startCode[i] = code
		code += count[i]
		code <<= 1
	}

	nextNode := uint16(1) // Node 0 is the root
	
	for i := range t.left {
		t.left[i] = 0
		t.right[i] = 0
	}

	for symbol := uint16(0); symbol < numSymbols; symbol++ {
		length := lengths[symbol]
		if length == 0 {
			continue 
		}

		currentCode := startCode[length]
		startCode[length]++

		node := uint16(0) 
		
		for bitPos := length; bitPos > 0; bitPos-- {
			bit := (currentCode >> (bitPos - 1)) & 1

			if bit == 0 {
				if t.left[node] == 0 {
					t.left[node] = nextNode
					nextNode++
				}
				node = t.left[node]
			} else {
				if t.right[node] == 0 {
					t.right[node] = nextNode
					nextNode++
				}
				node = t.right[node]
			}
		}

		t.left[node] = symbol
		t.right[node] = 0xFFFF // Marker for leaf node
	}

	return nil
}

// readSymbol reads bits from the stream and walks the tree to find the symbol.
func (t *huffmanTree) readSymbol(br *BitReader) (uint16, error) {
	node := uint16(0) // Start at root

	for t.right[node] != 0xFFFF { // While not a leaf
		bit, err := br.ReadBits(1)
		if err != nil {
			return 0, err
		}

		if bit == 0 {
			node = t.left[node]
		} else {
			node = t.right[node]
		}

		if node == 0 {
			return 0, errors.New("bad Huffman traversal: reached dead node")
		}
	}

	return t.left[node], nil
}

// ==========================================
// PART 3: THE DECODER ENGINE
// ==========================================

// lh5Decoder implements io.Reader to stream decompressed data.
type lh5Decoder struct {
	br        *BitReader
	remaining uint32
	blockSize uint16

	// LZSS Sliding Window (Ring Buffer)
	ringBuf [windowSize]byte
	ringPos int

	// Output buffering state
	outBuf []byte 

	// Huffman Trees attached to the decoder
	ncTree *huffmanTree // Literal/Length tree
	npTree *huffmanTree // Offset tree
}

// NewLH5Decoder initializes a new decoder for the -lh5- compression method.
func NewLH5Decoder(r io.Reader, originalSize uint32) io.Reader {
	var br io.ByteReader
	if b, ok := r.(io.ByteReader); ok {
		br = b
	} else {
		br = bufio.NewReader(r)
	}

	return &lh5Decoder{
		br:        &BitReader{r: br},
		remaining: originalSize,
		ncTree:    newHuffmanTree(510), // -lh5- uses max 510 symbols for literals/lengths
		npTree:    newHuffmanTree(40),  // -lh5- uses max 40 symbols for offsets
	}
}

func (d *lh5Decoder) readBlockHeader() error {
	// Read the size of the next block
	sizeBits, err := d.br.ReadBits(16)
	if err != nil {
		return err
	}
	d.blockSize = sizeBits

	// Read the Temporary Tree (NT), which is used to read the other trees
	// Read the Code Length Tree (NC)
	// Read the Offset Tree (NP)
	// (In a full implementation, this is where you read the arrays of bit-lengths 
	// from the file and call d.ncTree.buildTree() and d.npTree.buildTree())
	
	return d.readTrees() 
}

// readPtLen reads bit-lengths for a small tree directly from the stream.
// This closely translates the 'read_pt_len' C function.
func (d *lh5Decoder) readPtLen(nn uint16, nbit uint8, iSpecial uint16) (*huffmanTree, error) {
	numSymbols, err := d.br.ReadBits(nbit)
	if err != nil {
		return nil, err
	}

	lengths := make([]uint8, nn)

	if numSymbols == 0 {
		// Special Case: A tree with only one active symbol.
		singleSymbol, err := d.br.ReadBits(nbit)
		if err != nil {
			return nil, err
		}
		// In a single-node tree, we artificially set its length to 0.
		// The buildTree method will need to handle this edge case.
		if singleSymbol < nn {
			lengths[singleSymbol] = 0 // Wait for buildTree to handle
		}
	} else {
		i := uint16(0)
		for i < numSymbols && i < nn {
			// Read 3 bits directly. LHA sometimes stores bit lengths as unary/binary mixes,
			// but for PT/NP it reads 3 bits at a time.
			// Standard C implementation often looks like: 
			// c = GETBITS(3); if (c == 7) while (GETBIT()) c++;

			// Read the base 3 bits
			val, err := d.br.ReadBits(3)
			if err != nil {
				return nil, err
			}
			
			c := uint8(val)
			
			// If the value is exactly 7, switch to unary counting
			if c == 7 {
				for {
					bit, err := d.br.ReadBits(1)
					if err != nil { 
						return nil, err 
					}
					if bit == 0 { 
						break // Stop counting when we hit a 0
					}
					c++
				}
			}
			
			lengths[i] = c
			i++

			// Handle the "skip" code (iSpecial)
			if i == iSpecial {
				skipBits, err := d.br.ReadBits(2)
				if err != nil { return nil, err }
				skipCount := skipBits
				
				for j := uint16(0); j < skipCount && i < nn; j++ {
					lengths[i] = 0
					i++
				}
			}
		}
	}

	tree := newHuffmanTree(nn)
	err = tree.buildTree(lengths)
	return tree, err
}

// readCLen reads the lengths for the Literal/Length tree using the Pre-Tree.
func (d *lh5Decoder) readCLen(ptTree *huffmanTree) (*huffmanTree, error) {
	numNC, err := d.br.ReadBits(9) // Max 510 symbols
	if err != nil {
		return nil, err
	}

	lengths := make([]uint8, 510)

	if numNC == 0 {
		// Handle empty tree
		singleSymbol, err := d.br.ReadBits(9)
		if err != nil { return nil, err }
		if singleSymbol < 510 {
			lengths[singleSymbol] = 0
		}
	} else {
		i := uint16(0)
		for i < numNC && i < 510 {
			// Ask the Pre-Tree to decode the next symbol
			c, err := ptTree.readSymbol(d.br)
			if err != nil {
				return nil, err
			}

			if c <= 2 {
				// Codes 0, 1, and 2 are special "skip" commands in LHA
				var skipCount uint16
				if c == 0 {
					skipCount = 1
				} else if c == 1 {
					skipBits, err := d.br.ReadBits(4)
					if err != nil { return nil, err }
					skipCount = skipBits + 3
				} else if c == 2 {
					skipBits, err := d.br.ReadBits(9)
					if err != nil { return nil, err }
					skipCount = skipBits + 20
				}

				for j := uint16(0); j < skipCount && i < 510; j++ {
					lengths[i] = 0
					i++
				}
			} else {
				// It's an actual length value
				lengths[i] = uint8(c - 2)
				i++
			}
		}
	}

	tree := newHuffmanTree(510)
	err = tree.buildTree(lengths)
	return tree, err
}

// readTrees reads the dynamic Huffman tree definitions for the new block.
func (d *lh5Decoder) readTrees() error {
	// 1. Read the Pre-Tree (PT). 
	// Max 19 symbols, read using 5 bits for the count, 3 bits for special index.
	ptTree, err := d.readPtLen(19, 5, 3)
	if err != nil {
		return err
	}

	// 2. Decode the Literal/Length (NC) Tree using the Pre-Tree.
	// Max 510 symbols.
	ncTree, err := d.readCLen(ptTree)
	if err != nil {
		return err
	}
	d.ncTree = ncTree

	// 3. Read the Offset (NP) Tree.
	// In -lh5-, the NP tree lengths are read directly from the stream 
	// using the exact same logic as the Pre-Tree. Max 19 symbols.
	npTree, err := d.readPtLen(19, 5, 3)
	if err != nil {
		return err
	}
	d.npTree = npTree

	return nil
}

// Read implements the standard io.Reader interface.
func (d *lh5Decoder) Read(p []byte) (int, error) {
	if d.remaining == 0 && len(d.outBuf) == 0 {
		return 0, io.EOF
	}

	written := 0

	// 1. Flush any buffered output from a previous LZSS match
	if len(d.outBuf) > 0 {
		n := copy(p, d.outBuf)
		d.outBuf = d.outBuf[n:]
		written += n
		p = p[n:]
	}

	// 2. Main decompression loop
	for len(p) > 0 && d.remaining > 0 {
		
		// If we finished the previous block, load the next one
		if d.blockSize == 0 {
			if err := d.readBlockHeader(); err != nil {
				return written, err
			}
		}
		d.blockSize--

		// Read the next symbol
		symbol, err := d.ncTree.readSymbol(d.br)
		if err != nil {
			return written, err
		}

		if symbol < 256 {
			// Literal byte
			b := byte(symbol)
			
			d.ringBuf[d.ringPos] = b
			d.ringPos = (d.ringPos + 1) % windowSize
			
			p[0] = b
			p = p[1:]
			written++
			d.remaining--
		} else {
			// LZSS Match Length & Offset
			matchLength := int(symbol - 256 + 3) 
			
			offsetIndex, err := d.npTree.readSymbol(d.br)
			if err != nil {
				return written, err
			}
			
			matchOffset := int(offsetIndex)
			if offsetIndex > 1 {
				extraBits, err := d.br.ReadBits(uint8(offsetIndex - 1))
				if err != nil {
					return written, err
				}
				matchOffset = (1 << (offsetIndex - 1)) + int(extraBits)
			}

			// Copy the match from the dictionary
			for i := 0; i < matchLength; i++ {
				readPos := (d.ringPos - matchOffset - 1 + windowSize) % windowSize
				b := d.ringBuf[readPos]
				
				d.ringBuf[d.ringPos] = b
				d.ringPos = (d.ringPos + 1) % windowSize
				
				if len(p) > 0 {
					p[0] = b
					p = p[1:]
					written++
				} else {
					d.outBuf = append(d.outBuf, b)
				}
				d.remaining--
			}
		}
	}

	return written, nil
}
