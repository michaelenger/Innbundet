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

	"github.com/michaelenger/innbundet/log"
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

// Convert a source URL to its scheme and hostname
func getHostname(sourceUrl string) string {
	u, err := url.Parse(sourceUrl)
	if err != nil {
		return sourceUrl
	}

	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

// Ensure that the URL to an item is an absolute URL, adding in the base if not
func ensureAbsoluteUrl(baseUrl, itemUrl string) string {
	u, err := url.Parse(itemUrl)
	if u.Scheme != "" || err != nil {
		return itemUrl
	}

	// Use pure hostname if the item starts with "/"
	if strings.HasPrefix(itemUrl, "/") {
		baseUrl = getHostname(baseUrl)
	}

	itemUrl, _ = url.JoinPath(baseUrl, itemUrl)
	itemUrl, _ = url.PathUnescape(itemUrl)

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
		case html.StartTagToken, html.SelfClosingTagToken:
			tn, _ := tokenizer.TagName()
			if string(tn) != "link" {
				break
			}

			attrs := extractLinkElement(tokenizer)
			if attrs.Href == "" {
				break // no point in caring about this
			}

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

	// Haven't found any links, try the root URL
	if len(links) == 0 {
		siteUrl = getHostname(siteUrl)
		links, err = fetchLinkElements(siteUrl)
		if err != nil {
			return nil
		}
	}

	icon := ""
	currentSize := -1
	for _, link := range links {
		matched, err := regexp.MatchString("icon", link.Rel)
		if matched == false || err != nil {
			continue
		}

		size := 0
		sizes := strings.Split(link.Sizes, "x")
		if len(sizes) != 0 {
			size, _ = strconv.Atoi(sizes[0])
		}

		if size < currentSize {
			continue
		}

		if strings.HasPrefix(link.Href, "data:") {
			continue // don't want to deal with this
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

	log.Debug("Parsing URL: %s", siteUrl)
	log.Debug("Attempting to parse as feed...")
	_, error := parser.ParseURL(siteUrl)
	if error == nil {
		log.Debug("Successfully parsed as a feed")
		return []string{siteUrl}, nil // if we're able to parse it then it's a feed URL
	}

	log.Debug("Treating as website")
	log.Debug("Extracting link elements...")
	links, err := fetchLinkElements(siteUrl)
	if err != nil {
		return nil, err
	}
	log.Debug("Got %d link elements", len(links))

	var urls []string
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
		log.Debug("Added potential feed URL: %s", feedUrl)
		urls = append(urls, feedUrl)
	}

	return urls, nil
}
