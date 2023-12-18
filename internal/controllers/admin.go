package controllers

import (
	"bytes"
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/aidarkhanov/nanoid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

func Admin() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := models.GetDefaultSite("Admin Area")

		patients, err := models.GetAllPatients()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var rps []models.RemotePatient

		for _, patient := range patients {
			results, err := models.LoadPatientResults(patient.AuthId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			for _, result := range results {
				if result.Test == string(models.MMPI) {
					final, err := result.CalculateMMPI(patient.Name, result.Metric)
					if err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
					rps = append(rps, models.RemotePatient{ Id: patient.AuthId, Patient: patient.Name, Date: patient.Created.Format(time.DateOnly), Done: len(patient.Exams) <= 0, Test: models.MMPI, Indication: "", Score: -1, Max: -1, Results: final })
				} else {
					sc, err := strconv.Atoi(result.Metric)
						if err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
					test := models.Exam(result.Test)

					var ind string
					var mx int
					switch test {
					case models.ASQ:
						ind = models.CalcTestASQ(sc, patient.Name)
						mx = models.ASQ_MAX_SCORE
					case models.BAI:
						ind = models.CalcTestBAI(sc, patient.Name)
						mx = models.BAI_MAX_SCORE
					case models.BDI:
						ind = models.CalcTestBDI(sc, patient.Name)
						mx = models.BDI_MAX_SCORE
					case models.P3:
						ind = models.CalcTestP3(sc, patient.Name)
						mx = models.P3_MAX_SCORE
					}

					rps = append(rps, models.RemotePatient{ Id: patient.AuthId, Patient: patient.Name, Date: patient.Created.Format(time.DateOnly), Done: len(patient.Exams) <= 0, Test: models.Exam(result.Test), Indication: ind, Score: sc, Max: mx, Results: models.MMPIResults{} })
				}

			}
		}

		html, err := helpers.GeneratePage(views.Admin(data, rps))

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
			log.Warn("TODO: you need to implement this properly")
			log.Errorf("rendering index: %s", err)
		}

		return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())

	}

}
