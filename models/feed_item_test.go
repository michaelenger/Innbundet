package models

import (
	"testing"
)

func TestSimpleDescription(t *testing.T) {
	item := FeedItem{}
	tests := map[string]string{
		"This is a test.":                         "This is a test.",
		"  This is a test.    ":                   "This is a test.",
		"<p>This is a test.</p>":                  "This is a test.",
		"<p>This is <a href=\"\">a test</a>.</p>": "This is a test.",
		"<img src=\"\"> <p>This is a test.</p>":   "This is a test.",
	}

	for description, expected := range tests {
		item.Description = description
		result := item.SimpleDescription()
		if result != expected {
			t.Fatalf("received: \"%+v\" expected: \"%+v\"", result, expected)
		}
	}
}
