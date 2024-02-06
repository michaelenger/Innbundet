package models

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// An item in a feed
type FeedItem struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Link        string
	Description string
	Author      string
	Image       *string
	Published   time.Time
	FeedID      uint
	Feed        Feed
}

// Returns a simplified description, removing HTML tags and shortening it.
func (item FeedItem) SimpleDescription() string {
	var buf []byte
	description := bytes.NewBuffer(buf)

	tokenizer := html.NewTokenizer(strings.NewReader(item.Description))
loop:
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err != io.EOF {
				return fmt.Sprintf("Failed to parse description: %v", err)
			}

			break loop
		case html.TextToken:
			description.Write(tokenizer.Text())
		}
	}

	return strings.TrimSpace(description.String())
}
