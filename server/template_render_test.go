package server

import (
	"testing"
)

func TestTruncateString(t *testing.T) {
	tests := map[string]string{
		"This is a test.": "This is a...",
		"Test":            "Test",
		"Thisshouldnotbetruncatedduetomissingspaces": "Thisshouldnotbetruncatedduetomissingspaces",
		"Xactly ten": "Xactly ten",
	}

	for input, expected := range tests {
		result := truncateString(input, 10)
		if result != expected {
			t.Fatalf("received: \"%+v\" expected: \"%+v\"", result, expected)
		}
	}
}
