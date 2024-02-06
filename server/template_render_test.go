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
			t.Fatalf("testing \"%v\" received \"%v\" expected \"%v\"", input, result, expected)
		}
	}
}

func TestUrlHost(t *testing.T) {
	tests := map[string]string{
		"https://example.org":                     "example.org",
		"example.org":                             "example.org",
		"https://www.example.org":                 "example.org",
		"https://EXAMPLE.ORG/some/page?weird=yes": "example.org",
		"https://localhost/some/page?weird=yes":   "localhost",
	}

	for input, expected := range tests {
		result := urlHost(input)
		if result != expected {
			t.Fatalf("testing \"%v\" received \"%v\" expected \"%v\"", input, result, expected)
		}
	}
}
