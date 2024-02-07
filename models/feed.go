package models

import (
	"time"

	"gorm.io/gorm"
)

// A feed that is being subscribed to
type Feed struct {
	ID          uint `gorm:"primaryKey"`
	Url         string
	Title       string
	Link        string
	Description string
	Image       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Create a new feed or update if the feed already exists
func CreateOrUpdateFeed(db *gorm.DB, feed *Feed) (*Feed, bool, error) {
	var created bool
	var existingFeed *Feed
	var result *gorm.DB

	db.Where("url = ?", feed.Url).Limit(1).Find(&existingFeed)
	if existingFeed.ID != 0 {
		result = db.Model(existingFeed).Updates(Feed{
			Title:       feed.Title,
			Link:        feed.Link,
			Description: feed.Description,
			Image:       feed.Image,
		})
		feed = existingFeed
		created = false
	} else {
		result = db.Create(feed)
		created = true
	}

	return feed, created, result.Error
}
