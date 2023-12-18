package models

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type RemotePatient struct {
	Id         string
	Patient    string
	Date       string
	Done       bool
	Test       Exam
	Indication string
	Score      int
	Max        int
	Results    MMPIResults
}

type AdminResult struct {
	ID string
	Patient string
	Sex string
	Test string
	Metric string
	Duration int
	Created time.Time
	Aid string
}


func (results *AdminResult) Submit() error {
	statement := `INSERT INTO adminresults(id, patient, sex, test, answers, duration, created, aid) VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := db.Exec(statement, results.ID, results.Patient, results.Sex, results.Test, results.Metric, results.Duration, results.Aid)

	if err != nil {
		return err
	}

	return nil
}


func LoadAdminResult(aid string) ([]AdminResult, error) {
	statement := `SELECT * FROM adminresults WHERE aid=$1;`

	rows, err := db.Query(statement, aid)

	if err != nil {
		return []AdminResult{}, err
	}

	defer rows.Close()

	var results []AdminResult
	for rows.Next() {
		var result AdminResult
		err = rows.Scan(&result.ID, &result.Patient, &result.Sex, &result.Test, &result.Metric, &result.Duration, &result.Created, &result.Aid)
		if err != nil {
			return []AdminResult{}, err
		}

		results = append(results, result)

	}

	return results, nil
}

type Examination struct {
	ID int
	Sex string
	Test string
	Metric string
	Duration int
	Created time.Time
	Pid string
}

func (results *Examination) SubmitExamination() error {
	statement := `INSERT INTO examinations(sex, test, answers, duration, created, pid) VALUES($1, $2, $3, $4, $5, $6);`

	_, err := db.Exec(statement, results.Sex, results.Test, results.Metric, results.Duration, results.Created, results.Pid)

	if err != nil {
		return err
	}

	return nil
}


func LoadPatientResults(pid string) ([]Examination, error) {
	statement := `SELECT * FROM examinations WHERE pid=$1;`

	rows, err := db.Query(statement, pid)

	if err != nil {
		return []Examination{}, err
	}

	defer rows.Close()

	var results []Examination
	for rows.Next() {
		var result Examination
		err = rows.Scan(&result.ID, &result.Sex, &result.Test, &result.Metric, &result.Duration, &result.Created, &result.Pid)
		if err != nil {
			return []Examination{}, err
		}

		results = append(results, result)

	}

	return results, nil
}

func (ex *Examination) CalculateMMPI(patient string, answers string) (MMPIResults, error) {
	var results MMPIResults

	results.ID = ex.Pid
	results.Patient = patient
	results.Sex = ex.Sex
	results.Duration = ex.Duration

	filename := "data/scales.json"
	var scalesData MMPIScales

	scj, err := os.ReadFile(filename)
	if err != nil {
		return MMPIResults{}, err
	}

	err = json.Unmarshal(scj, &scalesData)
	if err != nil {
		return MMPIResults{}, err
	}

	var parsedAnswers []bool

	for _, answer := range strings.Split(answers, "") {
		if answer == "T" {
			parsedAnswers = append(parsedAnswers, true)
		} else {
			parsedAnswers = append(parsedAnswers, false)
		}
	}

	for _, category := range scalesData {

		results.Categories = append(results.Categories, MMPICategoryResult{
			Title:              category.Title,
			Scales:             make([]ScaleResult, 0),
			DerivedIndications: make([]string, 0),
		})

		for _, scale := range category.Items {
			err := grade(parsedAnswers, scale, ex.Sex, &results)
			if err != nil {
				return MMPIResults{}, err
			}
		}

		results.Categories[len(results.Categories)-1].deriveIndications(scalesData, &results)
	}

	return results, nil
}
