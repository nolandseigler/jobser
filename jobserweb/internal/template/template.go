package template

import (
	"embed"
	"fmt"
	htmpl "html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed templates
var tmplFS embed.FS

type Templates struct {
	templates map[string]*htmpl.Template
}

func New() *Templates {

	return &Templates{
		templates: map[string]*htmpl.Template{
			"dashboard": htmpl.Must(htmpl.ParseFS(tmplFS, "templates/dashboard.html", "templates/base.html")),
			"signup":    htmpl.Must(htmpl.ParseFS(tmplFS, "templates/signup.html", "templates/base.html")),
		},
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if tmpl, ok := t.templates[name]; ok {
		return tmpl.ExecuteTemplate(w, "base.html", data)
	}

	return fmt.Errorf("template not found; name: %s", name)
}
