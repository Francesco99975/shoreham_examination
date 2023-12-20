package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Locals() echo.HandlerFunc {
	return func(c echo.Context) error {

		data := models.GetDefaultSite("Local Patients")

		sess, err := session.Get("session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data for session")
		}

		results, err := models.LoadAdminResults(sess.Values["authid"].(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		var rps []models.RemotePatient

		for _, result := range results {

				if result.Test == string(models.MMPI) {
					final, err := result.LocalCalculateMMPI(result.Metric)
					if err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
					rps = append(rps, models.RemotePatient{ Id: result.ID, Patient: result.Patient, Date: result.Created.Format(time.DateOnly), Done: true, Test: models.Tests[4], Indication: "", Score: -1, Max: -1, Results: final })
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
						ind = models.CalcTestASQ(sc, result.Patient, result.Duration)
						mx = models.ASQ_MAX_SCORE
						ts = models.Tests[0]
					case models.BAI:
						ind = models.CalcTestBAI(sc, result.Patient, result.Duration)
						mx = models.BAI_MAX_SCORE
						ts = models.Tests[1]
					case models.BDI:
						ind = models.CalcTestBDI(sc, result.Patient, result.Duration)
						mx = models.BDI_MAX_SCORE
						ts = models.Tests[2]
					case models.P3:
						ind = models.CalcTestP3(sc, result.Patient, result.Duration)
						mx = models.P3_MAX_SCORE
						ts = models.Tests[3]
					}

					rps = append(rps, models.RemotePatient{ Id: result.ID, Patient: result.Patient, Date: result.Created.Format(time.DateOnly), Done: true, Test: ts, Indication: ind, Score: sc, Max: mx, Results: models.MMPIResults{} })
				}


		}

		html, err := helpers.GeneratePage(views.Locals(data, rps))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not parse page patient auth")
		}

		return c.Blob(200, "text/html; charset=utf-8", html)
	}
}
