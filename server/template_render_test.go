package server

import (
	"testing"
	"time"
)

func TestTimeAgo(t *testing.T) {
	tests := map[time.Time]string{
		time.Now().AddDate(-1, 0, 0):                 time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
		time.Now().AddDate(0, -1, 0):                 time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -7):                 time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		time.Now().AddDate(0, 0, -2):                 "2 days ago",
		time.Now().AddDate(0, 0, -1):                 "1 day ago",
		time.Now().Add(time.Hour * -12):              "12 hours ago",
		time.Now().Add(time.Hour * -1):               "1 hour ago",
		time.Now().Add(time.Minute * -40):            "40 minutes ago",
		time.Now().Add(time.Minute * -1):             "1 minute ago",
		time.Now().Add(-30):                          "now",
		time.Date(2112, 1, 1, 12, 0, 0, 0, time.UTC): "2112-01-01",
	}

	for input, expected := range tests {
		result := timeAgo(input)
		if result != expected {
			t.Fatalf("testing \"%v\" received \"%v\" expected \"%v\"", input, result, expected)
		}
	}
}

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
