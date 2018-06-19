package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main()  {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Something!\n")
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
