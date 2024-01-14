package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Health check endpoint handler
func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// Run the server
func RunServer() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/ping", ping)

	// Begin
	e.Logger.Fatal(e.Start(":8080"))
}
