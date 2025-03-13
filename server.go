package main

import (
	"html/template"
	"io"
	"net/http"
	"places/handlers"
	"github.com/labstack/echo/v4"
)

// TemplateRenderer helps Echo render templates with data
type TemplateRenderer struct {
	templates *template.Template
}

// Render implements echo.Renderer
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// Set up the custom renderer for Echo
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
    
	// Homepage Route
	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title": "jimmys.place",
		}
		return c.Render(http.StatusOK, "home", data)
	})

    handlers.RegisterBlogRoutes(e)

	// Static files (CSS, images, etc.)
	e.Static("/static", "static")

	// Start the server
	e.Logger.Fatal(e.Start(":2020"))
}

