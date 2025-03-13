package handlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"os"
)

// TemplateRenderer helps Echo render templates with data.
type TemplateRenderer struct {
	templates *template.Template
}

// Render implements echo.Renderer.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// RegisterBlogRoutes registers all blog-related routes.
func RegisterBlogRoutes(e *echo.Echo) {
	// Set up custom renderer for Echo.
	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	// Blog Homepage.
	e.GET("/blog", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title": "jimmys.place--blog",
		}
		return c.Render(http.StatusOK, "blog", data)
	})

	// Dynamic Blog Pages.
	e.GET("/blog/:page", func(c echo.Context) error {
		page := c.Param("page")
		filePath := fmt.Sprintf("templates/%s.html", page)

		// Check if the requested page exists.
		if !fileExists(filePath) {
			return c.String(http.StatusNotFound, "Page not found")
		}

		// Render the template with a dynamic title.
		data := map[string]interface{}{
			"Title": fmt.Sprintf("jimmys.place--%s", page),
		}
		return c.Render(http.StatusOK, page, data)
	})
}

// fileExists checks if a file exists and is not a directory.
func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

