package parser

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

// An HTML <link> element
type linkElement struct {
	Href  string
	Rel   string
	Sizes string
	Type  string
}

// Ensure that the URL to an item is an absolute URL, adding in the base if not
func ensureAbsoluteUrl(baseUrl, itemUrl string) string {
	u, err := url.Parse(itemUrl)
	if u.Scheme != "" || err != nil {
		return itemUrl
	}

	itemUrl, _ = url.JoinPath(baseUrl, itemUrl)

	return itemUrl
}

// Extract a <link> element from a tokenizer
func extractLinkElement(tokenizer *html.Tokenizer) linkElement {
	var attr linkElement

	for {
		key, value, more := tokenizer.TagAttr()

		switch string(key) {
		case "href":
			attr.Href = string(value)
		case "rel":
			attr.Rel = string(value)
		case "sizes":
			attr.Sizes = string(value)
		case "type":
			attr.Type = string(value)
		}

		if more == false {
			break
		}
	}

	return attr
}

// Fetch all the <link> elements in a website
func fetchLinkElements(url string) ([]linkElement, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Unexpected status code (%d)", resp.StatusCode))
	}

	var tags []linkElement

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

			attrs := extractLinkElement(tokenizer)
			tags = append(tags, attrs)
		}
	}

	return tags, nil
}

// Fetch the icon from a website, attempting to get the biggest one you can find
func fetchIcon(siteUrl string) *string {
	links, err := fetchLinkElements(siteUrl)
	if err != nil {
		return nil
	}

	icon := ""
	currentSize := 0
	for _, link := range links {
		matched, err := regexp.MatchString("icon", link.Rel)
		if matched == false || err != nil {
			continue
		}

		sizes := strings.Split(link.Sizes, "x")
		if len(sizes) == 0 {
			continue
		}

		size, err := strconv.Atoi(sizes[0])
		if size < currentSize || err != nil {
			continue
		}

		icon = link.Href
		currentSize = size
	}

	if icon == "" {
		return nil // found nothing :(
	}

	iconUrl := ensureAbsoluteUrl(siteUrl, icon)

	return &iconUrl
}

// Given a feed or web URL, find any available feeds.
func FindFeedUrls(siteUrl string) ([]string, error) {
	parser := gofeed.NewParser()

	_, error := parser.ParseURL(siteUrl)
	if error == nil {
		return []string{siteUrl}, nil // if we're able to parse it then it's a feed URL
	}

	var urls []string
	links, err := fetchLinkElements(siteUrl)
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		if link.Rel != "alternate" {
			continue
		}

		matched, err := regexp.MatchString("application/(rss|atom|json)", link.Type)
		if err != nil {
			return nil, err
		}
		if matched == false {
			continue
		}

		feedUrl := ensureAbsoluteUrl(siteUrl, link.Href)
		urls = append(urls, feedUrl)
	}

	return urls, nil
}
