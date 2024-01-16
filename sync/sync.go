package sync

import (
	"log"

	"github.com/michaelenger/innbundet/models"
	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

// Sync a given feed
func SyncFeed(db *gorm.DB, feed *models.Feed) error {
	logger := log.Default()
	parser := gofeed.NewParser()

	logger.Printf("Syncing %s (%d)", feed.Url, feed.ID)

	data, error := parser.ParseURL(feed.Url)
	if error != nil {
		return error
	}

	// Get/update feed items
	// TODO

	// Update feed metadata
	var image *string
	if data.Image != nil {
		image = &data.Image.URL
	}
	db.Model(feed).Updates(models.Feed{
		Title:       data.Title,
		Link:        data.Link,
		Description: data.Description,
		Image:       image,
	})
	logger.Print("..updated metadata")

	return nil
}
