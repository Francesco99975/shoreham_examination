package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type MMPIContent struct {
	Questions []string `json:"questions"`
}

func MMPI(admin bool) echo.HandlerFunc {

	data := models.Site{
		AppName:  "Shoreham Examination",
		Title:    "MMPI-2 Exam",
		Metatags: models.SEO{Description: "Examination tool", Keywords: "tools,exam"},
		Year:     time.Now().Year(),
	}

	filename := "data/mmpi2.json"
	var cnt *MMPIContent

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

	return func(c echo.Context) error {

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 1
		}

		if page < 1 || page > 23 {
			page = 1
		}

		questionsPerPage := 25
		startIndex := (page - 1) * questionsPerPage
		endIndex := startIndex + questionsPerPage

		if startIndex < 0 || startIndex >= len(cnt.Questions) {
			return echo.NewHTTPError(404, "Page not found", http.StatusNotFound)
		}

		if endIndex > len(cnt.Questions) {
			endIndex = len(cnt.Questions)
		}

		paginatedQuestions := cnt.Questions[startIndex:endIndex]

		buf := bytes.NewBuffer(nil)

		if page == 1 {
			var path string
			if admin {
				path = "/admin/mmpi"
			} else {
				path = "/examination/mmpi"
			}
			err = views.MMPI(data, admin, paginatedQuestions, page, path).Render(context.Background(), buf)

			if err != nil {
				return echo.NewHTTPError(404, "Page not found", http.StatusNotFound)
			}

			return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())
		} else {
			pid := c.QueryParam("pid")

			err = views.MMPIFormPartial(paginatedQuestions, page, pid).Render(context.Background(), buf)

			if err != nil {
				return echo.NewHTTPError(404, "Page not found", http.StatusNotFound)
			}

			return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())
		}
	}
}

func MMPICalc(admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var baseRedirectPath string
		if admin {
			baseRedirectPath = "/admin"
		} else {
			baseRedirectPath = "/examination"
		}
		answersPerPage := 25
		pid := c.FormValue("pid")
		patient := c.FormValue("patient")

		page, err := strconv.Atoi(c.FormValue("page"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data for page")
		}

		duration, err := strconv.Atoi(c.FormValue("duration"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data for duration")
		}

		sess, err := session.Get("session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data for session")
		}

		if len(patient) <= 0 && !admin {
			patient = sess.Values["patient"].(string)
		}

		var answers []string

		var limit int
		if page != 23 {
			limit = answersPerPage
		} else {
			limit = models.MMPI2_MAX_SCORE - (answersPerPage * 22)
		}

		// TESTING MMPI CODE ONLY

		startIndex := (page - 1) * limit
		endIndex := startIndex + limit

		test := strings.Split(models.MMPI_TEST_ANSWERS, "")[startIndex:endIndex]

		answers = append(answers, test...)

		// for i := 0; i < limit; i++ {
		// 	answer, err := strconv.Atoi(c.FormValue(fmt.Sprintf("%dA%d", page, i)))
		// 	if err != nil {
		// 		return echo.NewHTTPError(http.StatusBadRequest, "Insifficient data")
		// 	}

		// 	if answer > 1 || answer < 0 {
		// 		return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		// 	}

		// 	if answer == 0 {
		// 		answers = append(answers, "F")
		// 	} else {
		// 		answers = append(answers, "T")
		// 	}
		// }

		var newlocal models.LocalRes
		var newTemporal models.PatientRes

		if page == 1 {
			sex := c.FormValue("sex")
			if admin {
				newlocal = models.LocalRes{ID: uuid.NewV4().String(), Patient: patient, Sex: sex, Page: uint16(page + 1), Answers: strings.Join(answers, ""), Duration: duration, Aid: sess.Values["email"].(string)}
				err = newlocal.Save()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Could not save tempoary data: %s", err.Error()))
				}
				pid = newlocal.ID
			} else {
				newTemporal = models.PatientRes{Sex: sex, Page: uint16(page + 1), Answers: strings.Join(answers, ""), Duration: duration, Pid: sess.Values["authid"].(string)}
				err = newTemporal.PSave()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Could not save tempoary data for patient: %s", err.Error()))
				}
				pid = newTemporal.Pid
			}

			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/mmpi?page=%d&pid=%s", baseRedirectPath, page+1, pid))
		}

		if admin {
			newlocal, err = models.Load(pid)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not load tempoary data for local patient")
			}
			err = newlocal.Update(page+1, answers, duration)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Could not update tempoary data for local patient: %s", err.Error()))
			}
		} else {
			newTemporal, err = models.PLoad(sess.Values["authid"].(string))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not load temporary data for patient")
			}
			err = newTemporal.PUpdate(page+1, answers, duration)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Could not update tempoary data for patient: %s", err.Error()))
			}
		}

		if page == 23 {
			var results models.MMPIResults
			if admin {
				results, err = newlocal.Calculate()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during evaluation: %s", err.Error()))
				}
			} else {
				results, err = newTemporal.PCalculate(patient)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during evaluation for patient: %s", err.Error()))
				}
			}

			file, err := helpers.GeneratePDFMMPI(results)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during pdf generation: %s", err.Error()))
			}

			success, err := helpers.SendEmail("MMPI-2", "Attched a file with MMPI-2 results", results.Patient, file)

			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending: %s", err.Error()))
			}

			if success && admin {
				result := models.AdminResult { ID: newlocal.ID, Patient: patient, Sex: newlocal.Sex , Test: string(models.MMPI), Metric: strings.Join(answers, ""), Duration: duration, Created: time.Now(), Aid: sess.Values["email"].(string) }
				err = result.Submit()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Could not save admin test results")
				}
				return c.Redirect(http.StatusSeeOther, "/success")
			} else if success && !admin {
				pt, err := models.GetPatient(sess.Values["authid"].(string))
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error patient progression: %s", err.Error()))
				}
				exam, err := pt.NextExam()
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}

				result := models.Examination { Test: string(models.MMPI), Metric: strings.Join(answers, ""), Duration: duration, Created: time.Now(), Pid: sess.Values["authid"].(string) }
				err = result.SubmitExamination()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Could not save patient test results")
				}

				return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/examination?next=%s", exam))
			} else {
				return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error during email sending(failed): %s", err.Error()))
			}
		}

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s/mmpi?page=%d&pid=%s", baseRedirectPath, page+1, pid))
	}
}
