package controllers

import (
	"net/http"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		email := c.FormValue("email")
		password := c.FormValue("password")

		member, err := models.GetMemeber(email)

		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Member not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: wrong password")
		}

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not Authorized")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			// Secure: true, https
			// Domain: "",
			// SameSite: http.SameSiteDefaultMode,
		}

		sess.Values["authenticated"] = true
		sess.Save(c.Request(), c.Response())

		return c.JSON(http.StatusOK, "Authenticated")
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not Authorized")
		}

		sess.Values["authenticated"] = false

		return c.JSON(http.StatusOK, "Logged out")
	}
}
