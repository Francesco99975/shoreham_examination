package controllers

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/aidarkhanov/nanoid"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

func Admin() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Admin Area")

		html, err := helpers.GeneratePage(views.Admin(data))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page index")
		}


		return c.Blob(200, "text/html; charset=utf-8", html)
	}
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

		patient := models.Patient{AuthId: uuid.NewV4().String(), Name: name, Authcode: nanoid.New(), Exams: strings.Join(exams, ","), Created: time.Now()}

		err := models.CreatePatientExam(&patient)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create exam session")
		}

		buf := bytes.NewBuffer(nil)

		err = views.GenerationResults(patient.AuthId, patient.Authcode, patient.Name).Render(context.Background(), buf)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create view admin codes")
		}

		return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())

	}

}
