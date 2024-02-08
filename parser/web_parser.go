package parser

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

// An HTML link attribute
type linkAttributes struct {
	Rel  string
	Type string
	Href string
}

// Extract the attributes from a <link> element.
func extractLinkAttributes(tokenizer *html.Tokenizer) linkAttributes {
	var attr linkAttributes

	for {
		key, value, more := tokenizer.TagAttr()

		switch string(key) {
		case "href":
			attr.Href = string(value)
		case "rel":
			attr.Rel = string(value)
		case "type":
			attr.Type = string(value)
		}

		if more == false {
			break
		}
	}

	return attr
}

// Given a feed or web URL, find any available feeds.
func FindFeedUrls(url string) ([]string, error) {
	parser := gofeed.NewParser()

	_, error := parser.ParseURL(url)
	if error == nil {
		return []string{url}, nil // if we're able to parse it then it's a feed URL
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Unexpected status code (%d)", resp.StatusCode))
	}

	var urls []string

	tokenizer := html.NewTokenizer(resp.Body)
loop:
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				return nil, errors.New(fmt.Sprintf("Failed to parse HTML: %v", err))
			}

			break loop
		case html.StartTagToken:
			tn, _ := tokenizer.TagName()
			if string(tn) != "link" {
				break
			}

			attrs := extractLinkAttributes(tokenizer)
			if attrs.Rel != "alternate" {
				break
			}

			matched, err := regexp.MatchString("application/(rss|atom|json)", attrs.Type)
			if err != nil {
				return nil, err
			}
			if matched == false {
				break
			}

			urls = append(urls, attrs.Href)
		}
	}

	return urls, nil
}
