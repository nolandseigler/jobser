package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetDashboardHandler(c echo.Context) error {
	// TODO: remove me
	resp, err := http.Get("http://wordser:8080/api/v1/synonyms?word=friend")
	fmt.Printf("\n\n resp: %v, err: %v \n\n", resp, err)

	return c.Render(http.StatusOK, "dashboard", "")
}
