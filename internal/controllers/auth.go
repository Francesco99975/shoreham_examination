package controllers

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func PatientLogin() echo.HandlerFunc {
	return func(c echo.Context) error {

		authid := c.FormValue("authid")
		authcode := c.FormValue("authcode")

		patient, err := models.GetPatient(authid)

		if err != nil || patient.AuthId == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Patient not found")
		}

		if len(patient.Exams) <= 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Examination Already Concluded")
		}

		err = bcrypt.CompareHashAndPassword([]byte(patient.Authcode), []byte(authcode))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %s", err.Error()))
		}

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   true,
			Domain:   "shorehamex.dmz.urx.ink",
			SameSite: http.SameSiteDefaultMode,
		}

		sess.Values["authid"] = patient.AuthId
		sess.Values["patient"] = patient.Name
		sess.Values["examauth"] = true
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create auth session")
		}

		return c.Redirect(http.StatusSeeOther, "/examination")
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		email := c.FormValue("email")
		password := c.FormValue("password")

		member, err := models.GetMember(email)

		if err != nil || member.Email == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Member not found")
		}

		err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password))
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: wrong password")
		}

		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Server error on session")
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   true,
			Domain:   "shorehamex.dmz.urx.ink",
			SameSite: http.SameSiteDefaultMode,
		}

		sess.Values["email"] = member.Email
		sess.Values["authenticated"] = true
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create auth session")
		}

		return c.Redirect(http.StatusSeeOther, "/admin")
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Not Authorized")
		}

		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			Domain:   "shorehamex.dmz.urx.ink",
			SameSite: http.SameSiteDefaultMode,
		}

		sess.Values["email"] = ""
		sess.Values["authenticated"] = false
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create logout session")
		}

		return c.Redirect(http.StatusSeeOther, "/")
	}
}
