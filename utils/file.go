package utils

import (
	"os"
	"fmt"
	"path/filepath"
)

// isCreate: create directory if missing
func IsFileExists(path string, isCreate bool) (bool, error) {
	dir := filepath.Dir(path)
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if isCreate {
			err2 := os.MkdirAll(dir, 0750)
			if err2 != nil {
				return false, err2
			}
		} else {
			m := fmt.Errorf("No such directory: %s", dir)
			return false, m
		}
	} else if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("Not a directory: %s", dir)
	}
	info, err = os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
