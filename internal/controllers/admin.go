package controllers

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/aidarkhanov/nanoid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

func Admin() echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "Admin Area",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	return GeneratePage(views.Admin(data))
}

func GenerateCodes() echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("patient")

		availableExams := []string{"asq", "bai", "bdi", "p3", "mmpi"}
		var exams []string

		for i := 0; i < 5; i++ {
			asw := c.FormValue(availableExams[i])

			if asw == "on" {
				exams = append(exams, availableExams[i])
			}
		}

		if len(exams) <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Choose at least 1 test")
		}

		patient := models.Patient{AuthId: uuid.NewV4().String(), Name: name, Authcode: nanoid.New(), Exams: strings.Join(exams, ",")}

		models.CreatePatientExam(&patient)

		buf := bytes.NewBuffer(nil)

		err := views.GenerationResults(patient.AuthId, patient.Authcode, patient.Name).Render(context.Background(), buf)

		if err != nil {
			log.Warn("TODO: you need to implement this properly")
			log.Errorf("rendering index: %s", err)
		}

		return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())

	}

}
