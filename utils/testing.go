package utils

import (
	"fmt"
	"testing"
)

func Err[U any](msg [2]string, val U, e error) error {
	var m string
	if e == nil {
		m = fmt.Sprintf("%s %s", msg[0], msg[1])
	} else {
		m = fmt.Sprintf("%s (%v) %s", msg[0], e, msg[1])
	}
	return fmt.Errorf(m, val)
}

func TErr[U any](t *testing.T, msgPair [2]string, valPair [2]U) {
	m := fmt.Sprintf("%s %s", msgPair[0], msgPair[1])
	t.Errorf(m, valPair[0], valPair[1])
}
