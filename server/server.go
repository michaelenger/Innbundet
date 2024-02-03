package server

import (
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/models"
	"gorm.io/gorm"
)

// Custom server context
type ServerContext struct {
	echo.Context
	db     *gorm.DB
	config *config.Config
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
	result := ctx.db.Limit(ctx.config.ItemsPerPage).Offset(offset).Order("published desc").Find(&feedItems)
	if result.Error != nil {
		return result.Error
	}

	var total int64
	ctx.db.Model(&models.FeedItem{}).Count(&total)
	totalPages := int(math.Ceil(float64(total) / float64(ctx.config.ItemsPerPage)))

	return ctx.Render(http.StatusOK, "views/index.html", map[string]interface{}{
		"Config":     ctx.config,
		"FeedItems":  feedItems,
		"Page":       page,
		"TotalPages": totalPages,
	})
}

// Initialise the server
func Init(db *gorm.DB, conf *config.Config) (*echo.Echo, error) {
	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Custom context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ServerContext{c, db, conf}
			return next(cc)
		}
	})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	err := setupTemplateRenderer(e)
	if err != nil {
		return nil, err
	}

	// Routes
	e.GET("/", index)

	return e, nil
}
