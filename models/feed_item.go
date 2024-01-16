package models

import (
	"time"
)

// An item in a feed
type FeedItem struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Link        string
	Description *string
	Author      *string
	Published   time.Time
	FeedID      uint
	Feed        Feed
}
