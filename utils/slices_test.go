package utils

import (
	"fmt"
	"testing"
)

// Formatting patterns.
const FMT_TST_VAR = `Expected: "%v" Got: "%v"`

const TST_OK_SLICES = `Slices match.`
const TST_NO_SLICES_OK = `Slices match, but that's wrong.`

const TST_NO_SLICES = `Slices do not match.`
const TST_OK_SLICES_NO = `Slices do not match, but that's ok.`

func TestRemoveEmptyStrings_Must_Succeed(t *testing.T) {
	testStrings := []string{
		"", "test", "1234", "", "",
		"5678", "", "blah", "", "",
	}
	solution := []string{"test", "1234", "5678", "blah"}
	output := RemoveEmptyStrings(testStrings)
	if CompareSlices(output, solution) {
		t.Logf(TST_OK_SLICES)
	} else {
		m := fmt.Sprintf("%s %s", TST_NO_SLICES,
			FMT_TST_VAR)
		t.Errorf(m, output, solution)
	}
}

func TestRemoveEmptyStrings_Must_Fail(t *testing.T) {
}

func TestCompareSlices_Must_Succeed(t *testing.T) {
	as := []string{"abc", "123", "xyz"}
	ai := []int{1, -8, 43}

	bs := as
	bi := ai

	fs := []string{}
	fi := []int{1, 9, -4}

	passes := []bool{}
	passes = append(passes, CompareSlices(as, bs))
	passes = append(passes, CompareSlices(ai, bi))

	fails := []bool{}
	fails = append(fails, CompareSlices(as, fs))
	fails = append(fails, CompareSlices(ai, fi))

	for _, pass := range passes {
		if pass {
			t.Logf(TST_OK_SLICES)
		} else {
			t.Errorf(TST_NO_SLICES)
		}
	}

	for _, fail := range fails {
		if !fail {
			t.Logf(TST_NO_SLICES_OK)
		} else {
			t.Errorf(TST_OK_SLICES_NO)
		}
	}
}

func TestCompareSlices_Must_Fail(t *testing.T) {
}
