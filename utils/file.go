package utils

import (
	"os"
)

func IsFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		// doesn't exist
		return false, nil
	}
	// permissions issue
	return false, err
}
