package static

import (
	"embed"

	"github.com/labstack/echo/v4"
)

//go:embed static
var staticFS embed.FS

func RegisterStaticFS(e *echo.Echo) {
	e.StaticFS("/static", echo.MustSubFS(staticFS, "static"))
}
