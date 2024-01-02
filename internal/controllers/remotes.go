package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo/v4"
)

func Remotes() echo.HandlerFunc {
	return func(c echo.Context) error {

		data := models.GetDefaultSite("Remote Patients")

			patients, err := models.GetAllPatients()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var rps []models.RemotePatient

		for _, patient := range patients {
			var notDone []models.Exam
			if len(patient.Exams) > 0 {
				for _, test := range strings.Split(patient.Exams, ",") {
					notDone = append(notDone, models.Exam(test))
				}
			}
			results, err := models.LoadPatientResults(patient.AuthId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			for _, nd := range notDone {
					var mx int
					var ts models.TestSpecification
					switch nd {
					case models.ASQ:
						mx = models.ASQ_MAX_SCORE
						ts = models.Tests[0]
					case models.BAI:

						mx = models.BAI_MAX_SCORE
						ts = models.Tests[1]
					case models.BDI:

						mx = models.BDI_MAX_SCORE
						ts = models.Tests[2]
					case models.P3:
						mx = models.P3_MAX_SCORE
						ts = models.Tests[3]
					case models.MMPI:
						mx = -1
						ts = models.Tests[4]
					}

					rps = append(rps, models.RemotePatient{ Id: patient.AuthId, Patient: patient.Name, Date: patient.Created.Format(time.DateOnly), Done: false, Test: ts, Indication: "", Score: -1, Max: mx, Results: models.MMPIResults{} })
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

		html, err := helpers.GeneratePage(views.Remotes(data, rps))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page patient auth")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
