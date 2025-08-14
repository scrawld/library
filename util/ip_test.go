package util

import (
	"testing"
)

// TestParseIPPrefix
func TestParseIPPrefix(t *testing.T) {
	tests := []string{
		"192.168.1.5",
		"fe80::1",
		"2400:8901:e001:357:fa73:8e35:88b7:890b",
		"2a05:9400::631",
		"2a05:9400::632",
		"2a05:9403::7d1",
		"2a05:9403::7d2",
	}

	for _, ip := range tests {
		prefix, _ := ParseIPPrefix(ip)
		t.Logf("ParseIPPrefix IP: %s, Prefix: %s\n", ip, prefix)
	}
}
