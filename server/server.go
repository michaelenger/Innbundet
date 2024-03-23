package server

import (
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/log"
	"github.com/michaelenger/innbundet/models"
	"gorm.io/gorm"
)

// Custom server context
type ServerContext struct {
	echo.Context
	db     *gorm.DB
	config *config.Config
}

// Feed page - shows the entries in a single feed
func feed(c echo.Context) error {
	ctx := c.(*ServerContext)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get the feed
	feed := models.Feed{}
	result := ctx.db.First(&feed, id)
	if result.Error != nil {
		return result.Error
	}

	// Get the feed items
	feedItems := []models.FeedItem{}
	offset := (page - 1) * ctx.config.ItemsPerPage
	result = ctx.db.Preload("Feed").Where("feed_id = ?", id).Limit(ctx.config.ItemsPerPage).Offset(offset).Order("published desc").Find(&feedItems)
	if result.Error != nil {
		return result.Error
	}

	var total int64
	ctx.db.Model(&models.FeedItem{}).Where("feed_id = ?", id).Count(&total)
	totalPages := int(math.Ceil(float64(total) / float64(ctx.config.ItemsPerPage)))

	return renderTemplate(ctx, "views/feed.html", map[string]interface{}{
		"Feed":       feed,
		"FeedItems":  feedItems,
		"Page":       page,
		"TotalPages": totalPages,
	})
}

// Feeds page - shows a list of all the feeds
func feeds(c echo.Context) error {
	ctx := c.(*ServerContext)

	return renderTemplate(ctx, "views/feeds.html", map[string]interface{}{})
}

// Index page - show a list of recent feed items
func index(c echo.Context) error {
	ctx := c.(*ServerContext)

	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	feedItems := []models.FeedItem{}
	offset := (page - 1) * ctx.config.ItemsPerPage
	result := ctx.db.Preload("Feed").Limit(ctx.config.ItemsPerPage).Offset(offset).Order("published desc").Find(&feedItems)
	if result.Error != nil {
		return result.Error
	}

	var total int64
	ctx.db.Model(&models.FeedItem{}).Count(&total)
	totalPages := int(math.Ceil(float64(total) / float64(ctx.config.ItemsPerPage)))

	return renderTemplate(ctx, "views/index.html", map[string]interface{}{
		"FeedItems":  feedItems,
		"Page":       page,
		"TotalPages": totalPages,
	})
}

// Render a template in the given context ensuring the default data is present
func renderTemplate(ctx *ServerContext, template string, data map[string]interface{}) error {
	data["Config"] = ctx.config
	data["Context"] = ctx

	feeds := []models.Feed{}
	result := ctx.db.Order("title asc").Find(&feeds)
	if result.Error != nil {
		return result.Error
	}

	data["Feeds"] = feeds

	return ctx.Render(http.StatusOK, template, data)
}

// Initialise the server
func Init(db *gorm.DB, conf *config.Config) (*echo.Echo, error) {
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Custom context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ServerContext{c, db, conf}
			return next(cc)
		}
	})

	// Static files
	e.Static("/assets", "static")

	// Middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogError:    true,
		LogURI:      true,
		LogStatus:   true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				log.Info("Request: %s %s %d", v.Method, v.URI, v.Status)
			} else {
				log.Error("Request: %s %s %d: %v", v.Method, v.URI, v.Status, v.Error)
			}

			return nil
		},
	}))
	e.Use(middleware.Recover())
	err := setupTemplateRenderer(e)
	if err != nil {
		return nil, err
	}

	// Routes
	e.GET("/feeds/:id", feed)
	e.GET("/feeds", feeds)
	e.GET("/", index)

	return e, nil
}
