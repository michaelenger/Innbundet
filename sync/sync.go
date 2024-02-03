package sync

import (
	"log"
	"time"

	"github.com/michaelenger/innbundet/models"
	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

// Sync a given feed
func SyncFeed(db *gorm.DB, feed *models.Feed) error {
	var image *string

	logger := log.Default()
	parser := gofeed.NewParser()

	logger.Printf("Syncing %s (%d)", feed.Url, feed.ID)

	data, error := parser.ParseURL(feed.Url)
	if error != nil {
		return error
	}

	var feedAuthor string
	if len(data.Authors) != 0 {
		feedAuthor = data.Authors[0].Name
	}

	// Get/update feed items
	addCount := 0
	updateCount := 0
	for _, item := range data.Items {
		author := feedAuthor
		if len(item.Authors) != 0 {
			feedAuthor = item.Authors[0].Name
		}
		image = nil
		if item.Image != nil {
			image = &item.Image.URL
		}
		published := time.Now()
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		}

		var feedItem models.FeedItem
		db.Where("feed_id = ? AND link = ?", feed.ID, item.Link).Limit(1).Find(&feedItem)
		if feedItem.ID == 0 {
			db.Create(&models.FeedItem{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				Author:      author,
				Image:       image,
				Published:   published,
				Feed:        *feed,
			})
			addCount += 1
			continue
		}

		db.Model(&feedItem).Updates(models.FeedItem{
			Title:       item.Title,
			Description: item.Description,
			Author:      author,
			Image:       image,
			Published:   published,
		})
		updateCount += 1
	}

	if addCount != 0 {
		logger.Printf("..added %d feed items", addCount)
	}
	if updateCount != 0 {
		logger.Printf("..updated %d feed items", updateCount)
	}

	// Update feed metadata
	image = nil
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
