package controllers

import (
	"net/http"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

func Patient() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Patient")

		html, err := helpers.GeneratePage(views.Patient(data))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page patient auth")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
