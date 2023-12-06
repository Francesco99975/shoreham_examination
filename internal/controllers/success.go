package controllers

import (
	"time"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

func Success() echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "Staff",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	return GeneratePage(views.Success(data))
}
