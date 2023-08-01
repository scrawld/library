package util

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSliceShuffle(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7}
	SliceShuffle(s)
	t.Logf("s: %v", s)
}
