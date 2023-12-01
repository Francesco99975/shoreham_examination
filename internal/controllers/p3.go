package controllers

import (
	"time"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

func P3() echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "P3 Exam",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	admin := true
	questions := []string{"Is it it?", "Is it not?"}

	return GeneratePage(views.P3(data, admin, questions))
}
