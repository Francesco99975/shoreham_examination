package controllers

import (
	"net/http"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

func Staff() echo.HandlerFunc {

	return func(c echo.Context) error {
		data := models.GetDefaultSite("Staff")

		html, err := helpers.GeneratePage(views.Staff(data))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page staff auth")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
