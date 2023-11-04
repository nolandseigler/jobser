package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetSignupHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "signup", "")
}

type PostSignupRequest struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
}

func PostSignupHandler(auth Auther, db DBer) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(PostSignupRequest)
		if err := c.Bind(u); err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusBadRequest, "bad request")
		}
		passwordLen := len(u.Password)
		if passwordLen < 3 || passwordLen > 60 {
			return fmt.Errorf("password length must be >= 12 and <= 60; password length: %d", passwordLen)
		}
		_, err := db.CreateUserAccount(c.Request().Context(), u.Username, u.Password)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, "failed to create user account")
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
