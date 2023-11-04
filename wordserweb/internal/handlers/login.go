package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetLoginHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "login", "")
}

type PostLoginRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

func PostLoginHandler(auth Auther, db DBer) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(PostLoginRequest)
		if err := c.Bind(u); err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusBadRequest, "bad request")
		}

		jwt, err := auth.Login(c.Request().Context(), u.Username, u.Password)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, "failed to login")
		}

		// Set initial cookie. Auth middlewares in other endpoints keep it refreshed
		c.SetCookie(&http.Cookie{
			Name:    "session_token",
			Value:   jwt,
			Expires: time.Now().UTC().Add(60 * time.Minute),
		})

		return c.Redirect(http.StatusFound, "/dashboard")
	}
}
