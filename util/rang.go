package util

import (
	"math/rand"
	"time"
)

// RandNum
func RandNum(min, max int) int {
	if max < min || max == min {
		return min
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min+1)
}
