package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetDashboardHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard", "")
}
