package parser

import (
	"time"

	"github.com/michaelenger/innbundet/models"
	"github.com/mmcdole/gofeed"
)

// Extract a Feed from the feed data
func extractFeed(url string, data *gofeed.Feed) *models.Feed {
	var image *string

	link := data.Link
	if link == "" {
		link = getHostname(url)
	}

	if data.Image != nil {
		image = &data.Image.URL
	}

	if image == nil {
		image = fetchIcon(link)
	}

	feed := models.Feed{
		Title:       data.Title,
		Link:        link,
		Description: data.Description,
		Image:       image,
	}

	return &feed
}

// Extract FeedItems from the data
func extractFeedItems(data *gofeed.Feed) []*models.FeedItem {
	var feedAuthor string
	var items []*models.FeedItem

	if len(data.Authors) != 0 {
		feedAuthor = data.Authors[0].Name
	}

	for _, item := range data.Items {
		author := feedAuthor
		if len(item.Authors) != 0 {
			feedAuthor = item.Authors[0].Name
		}

		var image *string
		if item.Image != nil {
			image = &item.Image.URL
		} else if media, ok := item.Extensions["media"]; ok {
			if content, ok := media["content"]; ok && len(content) != 0 {
				if url, ok := content[0].Attrs["url"]; ok {
					image = &url
				}
			}
		}

		published := time.Now()
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		}

		feedItem := models.FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Author:      author,
			Image:       image,
			Published:   published,
		}

		items = append(items, &feedItem)
	}

	return items
}

// Parse the feed in a given URL, returning the appropriate Feed and FeedItems
func ParseFeed(url string) (*models.Feed, []*models.FeedItem, error) {
	parser := gofeed.NewParser()

	data, error := parser.ParseURL(url)
	if error != nil {
		return nil, nil, error
	}

	feed := extractFeed(url, data)
	feed.Url = url

	items := extractFeedItems(data)

	return feed, items, nil
}
