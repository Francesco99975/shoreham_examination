package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/helpers"
	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/Francesco99975/shorehamex/views"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

	qsj, err := helpers.ParseFile(filename)
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
			err = views.MMPI(data, admin, paginatedQuestions, page).Render(context.Background(), buf)

			if err != nil {
				log.Warn("TODO: you need to implement this properly")
				log.Errorf("rendering index: %s", err)
			}

			return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())
		} else {
			err = views.MMPIFormPartial(paginatedQuestions, page).Render(context.Background(), buf)

			if err != nil {
				log.Warn("TODO: you need to implement this properly")
				log.Errorf("rendering index: %s", err)
			}

			return c.Blob(200, "text/html; charset=utf-8", buf.Bytes())
		}
	}
}

func MMPICalc() echo.HandlerFunc {
	return func(c echo.Context) error {
		answersPerPage := 25
		patient := c.FormValue("patient")
		page, err := strconv.Atoi(c.FormValue("page"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		}

		var answers []string
		var limit int
		if page != 23 {
			limit = answersPerPage
		} else {
			limit = models.MMPI2_MAX_SCORE - (answersPerPage * 22)
		}

		for i := 0; i < limit; i++ {
			answer, err := strconv.Atoi(c.FormValue(fmt.Sprintf("%dA%d", page, i)))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Insifficient data")
			}

			if answer > 1 || answer < 0 {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
			}

			if answer == 0 {
				answers = append(answers, "F")
			} else {
				answers = append(answers, "T")
			}
		}

		sess, err := session.Get("session", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		}

		if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
		}

		var newlocal models.LocalRes

		if page == 1 {

			sex := c.FormValue("sex")
			newlocal = models.LocalRes{Patient: patient, Sex: sex, Page: uint16(page), Answers: strings.Join(answers, ""), Aid: sess.Values["email"].(string)}
			err = newlocal.Save()
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not save tempoary data")
			}

			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/mmpi?page=%d&patient=%s", page+1, patient))
		}

		newlocal, err = models.Load(patient)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not load tempoary data")
		}

		newlocal.Answers = strings.Join(append(strings.Split(newlocal.Answers, ""), answers...), "")

		if page == 23 {
			return c.Redirect(http.StatusSeeOther, "/")
		}

		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/mmpi?page=%d&patient=%s", page+1, patient))
	}
}
