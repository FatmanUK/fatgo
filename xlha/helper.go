package xlha

import (
	"fmt"
	"io"
	"strings"
)

const ERR_ARCH_SIZE_MISMATCH = `Size mismatch for %s: expected %d bytes, extracted %d`
const ERR_ARCH_404 = `End of archive reached. File not found.`

func ExtractFile(file io.Reader, n string, d *[]byte) error {
	lr := NewReader(file)
	for isDone := false; !isDone; {
		h, err := lr.Next()
		if err != nil {
			return err
		}
		*d, err = io.ReadAll(lr)
		if err != nil {
			return err
		}
		written := uint32(len(*d))
		oSz := h.OriginalSize
		if written != oSz {
			mFmt := ERR_ARCH_SIZE_MISMATCH
			return fmt.Errorf(mFmt, h.Name, oSz, written)
		}
		if (strings.Replace(n, "/", "\\", -1) == h.Name) {
			return nil
		}
	}
	return fmt.Errorf(ERR_ARCH_404)
}
