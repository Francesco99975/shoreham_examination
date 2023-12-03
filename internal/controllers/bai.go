package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

type BaiContent struct {
	Questions []string `json:"questions"`
}

func Bai(admin bool) echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "BAI Exam",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	filename := "data/bai.json"
	var cnt *BaiContent

	qsj, err := helpers.ParseFile(filename)
	if err != nil {
		fmt.Printf("error while reading json: %s", err.Error())
	}

	err = json.Unmarshal(qsj, &cnt)
	if err != nil {
		fmt.Printf("error while parsing json: %s", err.Error())
	}

	if len(cnt.Questions) <= 0 {
		return GeneratePage(views.ServerError(data, err))
	}

	return GeneratePage(views.Bai(data, admin, cnt.Questions))
}
