package controllers

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Examination() echo.HandlerFunc {
	return func(c echo.Context) error {
		exam := c.QueryParam("next")

		sess, err := session.Get("session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		}

		if auth, ok := sess.Values["examauth"].(bool); !ok || !auth {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data, not authenticated")
		}

		pt, err := models.GetPatient(sess.Values["authid"].(string))

		if err != nil || pt.AuthId == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Patient not found")
		}

		if exam == "cmp" {
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				// Secure: true, https
				// Domain: "",
				// SameSite: http.SameSiteDefaultMode,
			}

			sess.Values["authid"] = ""
			sess.Values["patient"] = ""
			sess.Values["examauth"] = false
			sess.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusSeeOther, "/success")
		}

		if len(exam) <= 0 {
			exam = pt.Peek()
		}

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/examination/%s", exam))
	}
}
