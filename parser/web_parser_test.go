package parser

import (
	"testing"
)

func TestEnsureAbsoluteUrl(t *testing.T) {
	baseUrl := "https://example.org"
	tests := map[string]string{
		"/test":               "https://example.org/test",
		"test/thing":          "https://example.org/test/thing",
		"http://yes.com/test": "http://yes.com/test",
	}

	for input, expected := range tests {
		result := ensureAbsoluteUrl(baseUrl, input)
		if result != expected {
			t.Fatalf("testing \"%v\" received \"%v\" expected \"%v\"", input, result, expected)
		}
	}

	baseUrl = "https://example.org/some/path"
	input := "/testing"
	expected := "https://example.org/testing"
	result := ensureAbsoluteUrl(baseUrl, input)
	if result != expected {
		t.Fatalf("testing \"%v\" received \"%v\" expected \"%v\"", input, result, expected)
	}
}
