package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
)

// Index page - show a list of recent feed items
func index(c echo.Context) error {
	db := db.DbManager()
	feedItems := []models.FeedItem{}

	result := db.Order("published desc").Find(&feedItems)

	if result.Error != nil {
		return result.Error
	}

	return c.Render(http.StatusOK, "views/index.html", map[string]interface{}{
		"feedItems": feedItems,
	})
}

// Initialise the server
func Init() (*echo.Echo, error) {
	// Echo instance
	e := echo.New()

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
