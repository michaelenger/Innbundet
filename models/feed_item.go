package models

import (
	"database/sql"
	"time"
)

// An item in a feed
type FeedItem struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Link        string
	Description sql.NullString
	Author      sql.NullString
	Published   time.Time
	FeedID      uint
	Feed        Feed
}
