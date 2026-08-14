package xlha_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"github.com/FatmanUK/fatgo/xlha"
)

// Test files:
// https://aminet.net/mods/inst/st-01.lha
// https://aminet.net/mods/inst/st-02.lha

// TestAminetArchives extracts real-world archives sequentially, ensuring
// the Huffman trees rebuild correctly and the CRC engine validates perfectly.
func TestAminetArchives(t *testing.T) {
	// The target test fixtures expected in your local workspace
	fixtures := []string{
		"st-01.lha",
		"st-02.lha",
	}

	for _, filename := range fixtures {
		t.Run(filename, func(t *testing.T) {
			path := filepath.Join("testdata", filename)

			// Open the local test fixture
			file, err := os.Open(path)
			if os.IsNotExist(err) {
				t.Skipf("Skipping test: copy %s into a 'testdata' folder to execute full verification", filename)
			}
			if err != nil {
				t.Fatalf("Failed to open fixture: %v", err)
			}
			defer file.Close()

			// Initialize your native extraction reader
			lhaReader := xlha.NewReader(file)

			fileCount := 0
			for {
				header, err := lhaReader.Next()
				if err == io.EOF {
					break // Successfully reached the end of the archive
				}
				if err != nil {
					t.Fatalf("[%s] failed at entry %d header parse: %v", filename, fileCount, err)
				}

				t.Logf("Processing file: %s (Method: %s, Original Size: %d)", header.Name, header.Method, header.OriginalSize)

				// Consume the uncompressed data stream entirely
				// This simultaneously updates the running table-driven CRC-16 engine
				written, err := io.Copy(io.Discard, lhaReader)
				if err != nil {
					t.Fatalf("Extraction failed for file %s: %v", header.Name, err)
				}

				if uint32(written) != header.OriginalSize {
					t.Errorf("Size mismatch for %s: expected %d bytes, extracted %d", header.Name, header.OriginalSize, written)
				}

				fileCount++
			}

			if fileCount == 0 {
				t.Errorf("Archive %s parsed without error but yielded 0 files.", filename)
			}

			t.Logf("Successfully verified %d files inside %s with zero CRC mismatches!", fileCount, filename)
		})
	}
}
