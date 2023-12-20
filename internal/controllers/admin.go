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
					rps = append(rps, models.RemotePatient{ Id: patient.AuthId, Patient: patient.Name, Date: patient.Created.Format(time.DateOnly), Done: len(patient.Exams) <= 0, Test: models.Tests[4], Indication: "", Score: -1, Max: -1, Results: final })
				} else {
					sc, err := strconv.Atoi(result.Metric)
						if err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
					test := models.Exam(result.Test)

					var ind string
					var mx int
					var ts models.TestSpecification
					switch test {
					case models.ASQ:
						ind = models.CalcTestASQ(sc, patient.Name, result.Duration)
						mx = models.ASQ_MAX_SCORE
						ts = models.Tests[0]
					case models.BAI:
						ind = models.CalcTestBAI(sc, patient.Name, result.Duration)
						mx = models.BAI_MAX_SCORE
						ts = models.Tests[1]
					case models.BDI:
						ind = models.CalcTestBDI(sc, patient.Name, result.Duration)
						mx = models.BDI_MAX_SCORE
						ts = models.Tests[2]
					case models.P3:
						ind = models.CalcTestP3(sc, patient.Name, result.Duration)
						mx = models.P3_MAX_SCORE
						ts = models.Tests[3]
					}

					rps = append(rps, models.RemotePatient{ Id: patient.AuthId, Patient: patient.Name, Date: patient.Created.Format(time.DateOnly), Done: len(patient.Exams) <= 0, Test: ts, Indication: ind, Score: sc, Max: mx, Results: models.MMPIResults{} })
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
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not create view admin codes")
		}

		return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())

	}

}
