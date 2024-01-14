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

// Run the server
func RunServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	err := setupTemplateRenderer(e)
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	// Routes
	e.GET("/", index)

	// Begin
	e.Logger.Fatal(e.Start(":8080"))
}
