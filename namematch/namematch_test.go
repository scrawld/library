package namematch_test

import (
	"testing"

	"test_kyc/library/namematch"
)

// TestNameSimilarity
func TestNameSimilarity(t *testing.T) {
	tests := []struct {
		a, b string
	}{
		{"John Smith", "John Smith"},
		{"Smith, John", "john smith"},
		{"Jonathan Doe", "Jonathon Doe"},
		{"John Smith", "Alice Brown"},
		{"John Smith", ""},
	}

	for _, tt := range tests {
		got := namematch.Similarity(tt.a, tt.b)
		t.Logf("%-15s %-15s %.2f", tt.a, tt.b, got)
	}
}
