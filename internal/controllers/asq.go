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

type AsqContent struct {
	Questions []string `json:"questions"`
	Multiq    []string `json:"multiq"`
}

func Asq(admin bool) echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "ASQ Exam",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	filename := "data/asq.json"
	var cnt *AsqContent

	qsj, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error while reading json: %s", err.Error())
	}

	err = json.Unmarshal(qsj, &cnt)
	if err != nil {
		fmt.Printf("error while parsing json: %s", err.Error())
	}

	if len(cnt.Multiq) <= 0 || len(cnt.Questions) <= 0 {
		return GeneratePage(views.ServerError(data, err))
	}

	return GeneratePage(views.Asq(data, admin, cnt.Questions, cnt.Multiq))
}

func AsqCalc() echo.HandlerFunc {
	return func(c echo.Context) error {
		patient := c.FormValue("patient")
		sex := c.FormValue("sex")

		var score int

		for i := 0; i < models.ASQ_MAX_SCORE; i++ {

			if i < 5 {
				asw, err := strconv.Atoi(c.FormValue(fmt.Sprintf("A%d", i)))
				if err != nil {
					echo.NewHTTPError(http.StatusBadRequest, "Bad request")
				}
				if asw == 1 {
					score++
				}
			} else {
				asw := c.FormValue(fmt.Sprintf("MA%d", i))

				if asw == "on" {
					score++
				}
			}
		}

		var indication string

		file, err := helpers.GeneratePDFGeneric("Anxiety Symptoms Questionnaire", patient, sex, indication, score)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during pdf generation: %s", err.Error()))
		}

		success, err := helpers.SendEmail("Anxiety Symptoms Questionnaire", patient, file)

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
