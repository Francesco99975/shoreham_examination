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
	uuid "github.com/satori/go.uuid"
)

type P3Content struct {
	Questions [][]string `json:"questions"`
}

func P3(admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
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
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = json.Unmarshal(qsj, &cnt)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if len(cnt.Questions) <= 0 {
			html, err := helpers.GeneratePage(views.ServerError(data, err))

			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page server error")
			}

			return c.Blob(200, "text/html; charset=utf-8", html)
		}

		var path string
		if admin {
			path = "/admin/p3"
		} else {
			path = "/examination/p3"
		}

		html, err := helpers.GeneratePage(views.P3(data, admin, cnt.Questions, path))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page P3 Exam")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
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
				return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
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

		indication := models.CompileBasicIndication(patient, percentage, "P3", duration, gravity)

		var id string
		sess, err := session.Get("session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data for session")
		}

		if admin {
			id = uuid.NewV4().String()
		} else {
			id = sess.Values["authid"].(string)
		}

		file, err := helpers.GeneratePDFGeneric("P3", id, patient, sex, duration, indication, score)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during pdf generation: %s", err.Error()))
		}

		success, err := helpers.SendEmail("P3", indication, patient, file)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending: %s", err.Error()))
		}

		if success && admin {
			result := models.AdminResult { ID: id, Patient: patient, Sex: sex, Test: string(models.P3), Metric: fmt.Sprint(score), Duration: duration, Created: time.Now(), Aid: sess.Values["email"].(string) }
			err = result.Submit()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not save admin test results")
			}
			return c.Redirect(http.StatusSeeOther, "/success")
		} else if success && !admin {
			pt, err := models.GetPatient(id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error patient progression: %s", err.Error()))
			}
			exam, err := pt.NextExam()
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			result := models.Examination { Sex: sex, Test: string(models.P3), Metric: fmt.Sprint(score), Duration: duration, Created: time.Now(), Pid: sess.Values["authid"].(string) }
			err = result.SubmitExamination()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not save patient test results")
			}

			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/examination?next=%s", exam))
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending(failed): %s", err.Error()))
		}
	}
}
