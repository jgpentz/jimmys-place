package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Route for the homepage
	e.GET("/", func(c echo.Context) error {
		return c.File("templates/index.html")
	})

	// Route for HTMX dynamic content
	e.GET("/blog", func(c echo.Context) error {
		return c.File("templates/blog.html")
	})

	// Route for HTMX dynamic content
	e.GET("/blog/page1", func(c echo.Context) error {
		return c.File("templates/page1.html")
	})

	// Route for HTMX dynamic content
	e.GET("/blog/page2", func(c echo.Context) error {
		return c.File("templates/page2.html")
	})

	// Route for HTMX dynamic content
	e.GET("/blog/page3", func(c echo.Context) error {
		return c.File("templates/page3.html")
	})

	e.Static("/static", "static")

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
