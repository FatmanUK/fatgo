package utils

import (
	"fmt"
	"strconv"
)

func StringFromAny(any any) string {
	return fmt.Sprintf("%v", any)
}

func Uint8FromHexString(s string) (uint8, error) {
	u, err := strconv.ParseUint(s, 16, 8)
	return uint8(u), err
}
