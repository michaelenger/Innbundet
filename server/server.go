package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Health check endpoint handler
func index(c echo.Context) error {
	return c.Render(http.StatusOK, "views/index.html", nil)
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
