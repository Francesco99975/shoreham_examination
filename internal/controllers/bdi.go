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

type BdiContent struct {
	Questions [][]string `json:"questions"`
}

func Bdi(admin bool) echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "BDI Exam",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	filename := "data/bdi.json"
	var cnt *BdiContent

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

	return GeneratePage(views.Bdi(data, admin, cnt.Questions))
}
