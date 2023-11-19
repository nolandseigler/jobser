package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO: Remove me. More hacks for homework
type SynonymsResp struct {
	Synonyms []string `json:"synonymns"`
}

func GetDashboardHandler(c echo.Context) error {
	// TODO: remove me
	resp, err := http.Get("http://wordser:8080/api/v1/synonyms?word=friend")
	fmt.Printf("\n\n resp: %v, err: %v \n\n", resp, err)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get synonyms from wordser; statusCode: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	syns := &SynonymsResp{}
	json.Unmarshal(data, syns)

	fmt.Printf("\n synonymns: %v \n", syns)

	return c.Render(http.StatusOK, "dashboard", "")
}
