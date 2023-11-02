package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetSignupHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", "")
}

type PostSignupRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

type PostSignupResponse struct {
	Username    string `json:"username" form:"username" query:"username"`
	AccessToken string `json:"access_token" form:"access_token" query:"access_token"`
}

func PostSignupHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(PostSignupRequest)
		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		username := "TODO: From DB"
		accessToken := "TODO: Mint JWT"

		return c.JSON(
			http.StatusCreated,
			PostSignupResponse{
				Username:    username,
				AccessToken: accessToken,
			},
		)
	}
}
