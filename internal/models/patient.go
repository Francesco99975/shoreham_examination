package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Patient struct {
	AuthId   string
	Name     string
	Authcode string
	Exams    string
	Created  time.Time
}

func CreatePatientExam(patient *Patient) error {

	statement := `INSERT INTO patients(authid, name, authcode, exams, created) VALUES($1, $2, $3, $4, $5);`

	hashedAuthCode, err := bcrypt.GenerateFromPassword([]byte(patient.Authcode), 12)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(statement, patient.AuthId, patient.Name, hashedAuthCode, patient.Exams, patient.Created)

	return err
}

func GetPatient(authid string) (Patient, error) {
	var patient Patient

	statement := `SELECT * FROM patients WHERE authid=$1;`

	rows, err := db.Query(statement, authid)

	if err != nil {
		return Patient{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&patient.AuthId, &patient.Name, &patient.Authcode, &patient.Exams, &patient.Created)
		if err != nil {
			return Patient{}, err
		}

	}

	return patient, nil
}

func (patient *Patient) Peek() string {
	return strings.Split(patient.Exams, ",")[0]
}

func (patient *Patient) NextExam() (string, error) {
	statement := `UPDATE patients SET exams=$1 WHERE authid=$2;`

	queued := strings.Split(patient.Exams, ",")

	if len(queued) <= 1 {
		patient.Exams = ""
		_, err := db.Exec(statement, patient.Exams, patient.AuthId)
		if err != nil {
			return "", err
		}
		return "cmp", nil
	}

	patient.Exams = strings.Join(queued[1:], ",")

	exam := strings.Split(patient.Exams, ",")[0]

	_, err := db.Exec(statement, patient.Exams, patient.AuthId)

	if err != nil {
		return "", err
	}

	return exam, nil
}

type PatientRes struct {
	ID       int
	Sex      string
	Page     uint16
	Answers  string
	Duration int
	Pid      string
}

func (temporal *PatientRes) PSave() error {
	statement := `INSERT INTO patientres(sex, page, answers, duration, pid) VALUES($1, $2, $3, $4, $5);`

	_, err := db.Exec(statement, temporal.Sex, temporal.Page, temporal.Answers, temporal.Duration, temporal.Pid)

	if err != nil {
		return err
	}

	return nil
}

func (temporal *PatientRes) PUpdate(newPage int, newAnswers []string, duration int) error {
	statement := `UPDATE patientres SET page=$1,answers=$2,duration=$3 WHERE pid=$4;`

	temporal.Page = uint16(newPage)
	temporal.Answers = strings.Join(append(strings.Split(temporal.Answers, ""), newAnswers...), "")
	temporal.Duration += duration

	_, err := db.Exec(statement, temporal.Page, temporal.Answers, temporal.Duration, temporal.Pid)

	if err != nil {
		return err
	}

	return nil
}

func (temporal *PatientRes) PCalculate(patient string) (MMPIResults, error) {
	var results MMPIResults

	results.ID = temporal.Pid
	results.Patient = patient
	results.Sex = temporal.Sex
	results.Duration = temporal.Duration

	filename := "data/scales.json"
	var scalesData MMPIScales

	scj, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error while reading json: %s", err.Error())
	}

	err = json.Unmarshal(scj, &scalesData)
	if err != nil {
		fmt.Printf("error while parsing json: %s", err.Error())
	}

	var parsedAnswers []bool

	for _, answer := range strings.Split(temporal.Answers, "") {
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
			err := grade(parsedAnswers, scale, temporal.Sex, &results)
			if err != nil {
				return MMPIResults{}, err
			}
		}

		results.Categories[len(results.Categories)-1].deriveIndications(scalesData, &results)
	}

	return results, nil
}

func PLoad(id string) (PatientRes, error) {
	statement := `SELECT * FROM patientres WHERE pid=$1;`

	rows, err := db.Query(statement, id)

	if err != nil {
		return PatientRes{}, err
	}

	defer rows.Close()

	var temporal PatientRes
	for rows.Next() {
		err = rows.Scan(&temporal.ID, &temporal.Sex, &temporal.Page, &temporal.Answers, &temporal.Duration, &temporal.Pid)
		if err != nil {
			return PatientRes{}, err
		}

	}

	return temporal, nil
}
