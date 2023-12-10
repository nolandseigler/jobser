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
			"login":     htmpl.Must(htmpl.ParseFS(tmplFS, "templates/login.html", "templates/base.html")),
			"analysis":  htmpl.Must(htmpl.ParseFS(tmplFS, "templates/analysis.html")),
		},
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if tmpl, ok := t.templates[name]; ok {
		if name == "analysis" {
			return tmpl.ExecuteTemplate(w, "analysis.html", data)
		}
		return tmpl.ExecuteTemplate(w, "base.html", data)
	}

	return fmt.Errorf("template not found; name: %s", name)
}
