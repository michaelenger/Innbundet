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
		"/feed?type=rss":      "https://example.org/feed?type=rss",
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

func TestGetHostname(t *testing.T) {
	tests := map[string]string{
		"https://example.org/":                  "https://example.org",
		"https://example.org/test":              "https://example.org",
		"http://yes.com/test/of/a/longer?thing": "http://yes.com",
		"ftp://localhost:420":                   "ftp://localhost:420",
	}

	for input, expected := range tests {
		result := getHostname(input)
		if result != expected {
			t.Fatalf("testing \"%v\" received \"%v\" expected \"%v\"", input, result, expected)
		}
	}
}
