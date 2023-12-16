package models

import "time"

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
	Test string
	Metric string
	Duration int
	Created time.Time
	Pid string
}

func (results *Examination) SubmitExamination() error {
	statement := `INSERT INTO examinations(test, answers, duration, created, pid) VALUES($1, $2, $3, $4, $5);`

	_, err := db.Exec(statement, results.Test, results.Metric, results.Duration, results.Created, results.Pid)

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
		err = rows.Scan(&result.ID, &result.Test, &result.Metric, &result.Duration, &result.Created, &result.Pid)
		if err != nil {
			return []Examination{}, err
		}

		results = append(results, result)

	}

	return results, nil
}
