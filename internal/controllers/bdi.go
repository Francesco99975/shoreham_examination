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
	"github.com/labstack/echo-contrib/session"
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

func BdiCalc(admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		patient := c.FormValue("patient")
		sex := c.FormValue("sex")
		duration, err := strconv.Atoi(c.FormValue("duration"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		}

		if len(patient) <= 0 {
			sess, err := session.Get("session", c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
			}
			patient = sess.Values["patient"].(string)
		}

		var score int

		for i := 0; i < 11; i++ {

			asw, err := strconv.Atoi(c.FormValue(fmt.Sprintf("A%d", i)))
			if err != nil {
				echo.NewHTTPError(http.StatusBadRequest, "Bad request")
			}

			score += asw

		}

		rawPercentage := float64(score) / float64(models.ASQ_MAX_SCORE) * 100.0

		var gravity string
		if rawPercentage <= 9 {
			gravity = "normal"
		} else if rawPercentage >= 10 && rawPercentage <= 18 {
			gravity = "moderate"
		} else {
			gravity = "severe"
		}

		percentage := fmt.Sprintf("%.2f", rawPercentage) + "%"

		indication := models.CompileBasicIndication(patient, percentage, "Beck Depression Inventory", gravity)

		file, err := helpers.GeneratePDFGeneric("Beck Depression Inventory", patient, sex, duration, indication, score)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during pdf generation: %s", err.Error()))
		}

		success, err := helpers.SendEmail("Beck Depression Inventory", patient, file)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending: %s", err.Error()))
		}
		if success && admin {
			return c.Redirect(http.StatusSeeOther, "/success")
		} else if success && !admin {
			return c.Redirect(http.StatusSeeOther, "/examination")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending(failed): %s", err.Error()))
		}
	}
}
