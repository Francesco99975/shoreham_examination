package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

type P3Content struct {
	Questions [][]string `json:"questions"`
}

func P3(admin bool) echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "P3 Exam",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	filename := "data/p3.json"
	var cnt *P3Content

	qsj, err := os.ReadFile(filename)
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

	return GeneratePage(views.P3(data, admin, cnt.Questions))
}
