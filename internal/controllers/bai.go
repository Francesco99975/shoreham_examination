package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

	return GeneratePage(views.Bai(data, admin, cnt.Questions))
}

func BaiCalc() echo.HandlerFunc {
	return func(c echo.Context) error {
		patient := c.FormValue("patient")
		sex := c.FormValue("sex")
		duration, err := strconv.Atoi(c.FormValue("duration"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		}

		var score int

		for i := 0; i < 21; i++ {

			asw, err := strconv.Atoi(c.FormValue(fmt.Sprintf("A%d", i)))
			if err != nil {
				echo.NewHTTPError(http.StatusBadRequest, "Bad request")
			}

			score += asw

		}

		var indication string

		file, err := helpers.GeneratePDFGeneric("Beck Anxiety Inventory", patient, sex, duration, indication, score)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during pdf generation: %s", err.Error()))
		}

		success, err := helpers.SendEmail("Beck Anxiety Inventory", patient, file)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending: %s", err.Error()))
		}

		if success {
			return c.Redirect(http.StatusSeeOther, "/success")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending(failed): %s", err.Error()))
		}
	}
}
