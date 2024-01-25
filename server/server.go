package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaelenger/innbundet/models"
	"gorm.io/gorm"
)

// Custom server context
type ServerContext struct {
	echo.Context
	db *gorm.DB
}

// Index page - show a list of recent feed items
func index(c echo.Context) error {
	ctx := c.(*ServerContext)

	feedItems := []models.FeedItem{}

	result := ctx.db.Order("published desc").Find(&feedItems)

	if result.Error != nil {
		return result.Error
	}

	return ctx.Render(http.StatusOK, "views/index.html", map[string]interface{}{
		"feedItems": feedItems,
	})
}

// Initialise the server
func Init(db *gorm.DB) (*echo.Echo, error) {
	// Echo instance
	e := echo.New()

	// Custom context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &ServerContext{c, db}
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
