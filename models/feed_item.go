package models

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/net/html"
	"gorm.io/gorm"
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

// Create a new feed item or update if the item already exists
func CreateOrUpdateFeedItem(db *gorm.DB, item *FeedItem) (*FeedItem, error) {
	var existingItem *FeedItem
	var result *gorm.DB

	db.Where("feed_id = ? AND link = ?", item.Feed.ID, item.Link).Limit(1).Find(&existingItem)
	if existingItem.ID != 0 {
		result = db.Model(existingItem).Updates(FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Author:      item.Author,
			Image:       item.Image,
			Published:   item.Published,
		})
		item = existingItem
	} else {
		result = db.Create(item)
	}

	return item, result.Error
}
