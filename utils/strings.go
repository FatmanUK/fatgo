package utils

import (
	"fmt"
	"strconv"
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io"
)

func StringFromAny(any any) string {
	return fmt.Sprintf("%v", any)
}

func Uint8FromHexString(s string) (uint8, error) {
	u, err := strconv.ParseUint(s, 16, 8)
	return uint8(u), err
}

// Compress shortens a string using zlib and encodes it to Base64
func Compress(src string) (string, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write([]byte(src))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	// Use RawURLEncoding to remove padding and make it URL safe
	return base64.RawURLEncoding.EncodeToString(b.Bytes()), nil
}

func Decompress(src string) (string, error) {
	cBytes, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	b := bytes.NewReader(cBytes)
	r, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}
	defer r.Close()
	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
