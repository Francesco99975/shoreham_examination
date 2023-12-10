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

	var path string
	if admin {
		path = "/admin/p3"
	} else {
		path = "/examination/p3"
	}

	return GeneratePage(views.P3(data, admin, cnt.Questions, path))
}

func P3Calc(admin bool) echo.HandlerFunc {
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

		for i := 0; i < 44; i++ {

			asw, err := strconv.Atoi(c.FormValue(fmt.Sprintf("A%d", i)))
			if err != nil {
				echo.NewHTTPError(http.StatusBadRequest, "Bad request")
			}

			score += asw

		}

		score = score - models.P3_ADJUST

		rawPercentage := float64(score) / float64(models.P3_MAX_SCORE) * 100.0

		var gravity string
		if rawPercentage <= 30 {
			gravity = "normal"
		} else if rawPercentage >= 31 && rawPercentage <= 50 {
			gravity = "moderate"
		} else {
			gravity = "severe"
		}

		percentage := fmt.Sprintf("%.2f", rawPercentage) + "%"

		indication := models.CompileBasicIndication(patient, percentage, "P3", gravity)

		file, err := helpers.GeneratePDFGeneric("P3", patient, sex, duration, indication, score)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during pdf generation: %s", err.Error()))
		}

		success, err := helpers.SendEmail("P3", indication, patient, file)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending: %s", err.Error()))
		}

		if success && admin {
			return c.Redirect(http.StatusSeeOther, "/success")
		} else if success && !admin {
			sess, err := session.Get("session", c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid data for session")
			}
			pt, err := models.GetPatient(sess.Values["authid"].(string))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error patient progression: %s", err.Error()))
			}
			exam, err := pt.NextExam()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/examination?next=%s", exam))
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending(failed): %s", err.Error()))
		}
	}
}
