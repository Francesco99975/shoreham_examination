package models

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
)

type Exam string

const (
	ASQ Exam = "asq"
	BAI Exam = "bai"
	BDI Exam = "bdi"
	P3  Exam = "p3"
)

const ASQ_MAX_SCORE int = 38
const BAI_MAX_SCORE int = 57
const BDI_MAX_SCORE int = 30
const P3_MAX_SCORE int = 85
const P3_ADJUST int = 44 // To Be subtracted from p3 score
const MMPI2_MAX_SCORE int = 567

type BasicExamResults struct {
	Score       int32
	Indications string
	Examinator  Exam
}

const MMPI_TEST_ANSWERS string = "TTFTFFTFFFFFFTFTFFFTTTFFFTFFTFTFFTFTTTTFTFFFFTTTFTTTTFTTTFTFFTTFTTTFTTTFTFF" +
	"TTFFFTTFFTFTTTFTTTFFFFTFFFFFFTFTTFTFFFFTFTTFTFFTFFFTTTTFFFFTFTFTFFFTFFFTFFT" +
	"FFFTTTFTTFTFTTFTTFFFTTFFTTFTTTTFFFTTFFFTFTFFTTFFTFFFFTFFTFTFFFFTTFTFTFFTFTF" +
	"TFFTFTTTFTTFTFFFFTFFFFFTTTFTFFTFTTTFFTTTTTTTFFFTFTTTTTFFFFTTFFFTFFFTFTFFFFF" +
	"TFTFFTFTTTTFFTFTFFFTFFFFTTFFFTFFFFFFTFTFTFFTFFTTTFFFTFFFTFFFFTFTTFFTFFTTFTF" +
	"TTTFTFFTTTTTFFTTFFTFFFFTTTFFTTFFTTFTTFFFFFTFTFFTTFFTFFFTTTTFFTFTTFTFFFTFFFT" +
	"TTFTTFFTTFTTTTTFTFTFFTTTTTFFTTFTTFTTTFTFTTFFFFTFTFTTFTTTFFTFTFFTTTTTTTFFTFF" +
	"TTFFTFTTFFFFFTFTFFFFFTFFTFFTTFFFTFFFTFTTTF"

type ScaleResult struct {
	ScaleName        string
	ScaleDescription string
	ScalePupose      string
	Score            int32
}

type MMPICategoryResult struct {
	Title              string
	Scales             []ScaleResult
	DerivedIndications []string
}

type MMPIResults struct {
	Patient    string
	Sex        string
	Categories []MMPICategoryResult
	Duration   int
}

type Indications struct {
	Four5__Mf__64                 []string `json:"45<=Mf<=64"`
	Five5__Hs__64                 []string `json:"55<=Hs<=64"`
	Five5__Hs__74                 []string `json:"55<=Hs<=74"`
	Five5__Hy__64                 []string `json:"55<=Hy<=64"`
	Five5__Pa__64                 []string `json:"55<=Pa<=64"`
	Five5__Pd__64                 []string `json:"55<=Pd<=64"`
	Five5__Pt__64                 []string `json:"55<=Pt<=64"`
	Five5__Sc__64                 []string `json:"55<=Sc<=64"`
	Six5__D__74                   []string `json:"65<=D<=74"`
	Six5__Hs__74                  []string `json:"65<=Hs<=74"`
	Six5__Hy__74                  []string `json:"65<=Hy<=74"`
	Six5__Pa__74                  []string `json:"65<=Pa<=74"`
	Six5__Pd__74                  []string `json:"65<=Pd<=74"`
	Six5__Pt__74                  []string `json:"65<=Pt<=74"`
	Six5__Sc__74                  []string `json:"65<=Sc<=74"`
	D__75                         []string `json:"D>=75"`
	Fb_F_20                       []string `json:"Fb>F+20"`
	Fp__100____VRIN_70____TRIN_70 string   `json:"Fp>=100 && VRIN<70 && TRIN<70"`
	Hs__75                        []string `json:"Hs>=75"`
	Hy__75                        []string `json:"Hy>=75"`
	Ma__55__64                    []string `json:"Ma<=55<=64"`
	Ma__65__74                    []string `json:"Ma<=65<=74"`
	Ma__75                        []string `json:"Ma>=75"`
	Mf_45                         []string `json:"Mf<45"`
	Mf__65                        []string `json:"Mf>=65"`
	Pa__75                        []string `json:"Pa>=75"`
	Pd__75                        []string `json:"Pd>=75"`
	Pt__75                        []string `json:"Pt>=75"`
	Sc__75                        []string `json:"Sc>=75"`
	Si_45                         []string `json:"Si<45"`
	Si__55__64                    []string `json:"Si<=55<=64"`
	Si__65__74                    []string `json:"Si<=65<=74"`
	Si__75                        []string `json:"Si>=75"`
	TRIN___80____TRIN__80         []string `json:"TRIN<=-80 || TRIN>=80"`
	VRIN_40                       []string `json:"VRIN<40"`
	VRIN__80                      []string `json:"VRIN>=80"`
}

type Scale struct {
	Answers      [][]interface{} `json:"answers"`
	BaseScore    int32           `json:"baseScore"`
	Code         string          `json:"code"`
	Comment      string          `json:"comment"`
	Gender       string          `json:"gender"`
	Indications  Indications     `json:"indications"`
	KCorrection  float32         `json:"kCorrection"`
	Name         string          `json:"name"`
	ScoreOffsets struct {
		Female int32 `json:"female"`
		Male   int32 `json:"male"`
	} `json:"scoreOffsets"`
	SubScales []Scale `json:"subScales"`
	TScores   struct {
		Female []int32 `json:"female"`
		Male   []int32 `json:"male"`
	} `json:"tScores"`
	Text  string `json:"text"`
	Title string `json:"title"`
}

type MMPIScales []struct {
	Items []Scale `json:"items"`
	Title string  `json:"title"`
}

type LocalRes struct {
	Patient  string
	Sex      string
	Page     uint16
	Answers  string
	Duration int
	Aid      string
}

func (local *LocalRes) Save() error {
	statement := `INSERT INTO localres(patient, sex, page, answers, duration, aid) VALUES($1, $2, $3, $4, $5, $6);`

	_, err := db.Exec(statement, local.Patient, local.Sex, local.Page, local.Answers, local.Duration, local.Aid)

	if err != nil {
		return err
	}

	return nil
}

func (local *LocalRes) Update(newPage int, newAnswers []string, duration int) error {
	statement := `UPDATE localres SET page=$1,answers=$2,duration=$3 WHERE patient=$4;`

	local.Page = uint16(newPage)
	local.Answers = strings.Join(append(strings.Split(local.Answers, ""), newAnswers...), "")
	local.Duration += duration

	_, err := db.Exec(statement, local.Page, local.Answers, local.Duration, local.Patient)

	if err != nil {
		return err
	}

	return nil
}

func (local *LocalRes) Calculate() (MMPIResults, error) {
	var results MMPIResults

	results.Patient = local.Patient
	results.Sex = local.Sex
	results.Duration = local.Duration

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

	for _, answer := range strings.Split(local.Answers, "") {
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
			err := grade(parsedAnswers, scale, local.Sex, &results)
			if err != nil {
				return MMPIResults{}, err
			}
		}

		results.Categories[len(results.Categories)-1].deriveIndications(scalesData, &results)
	}

	return results, nil
}

func Load(patient string) (LocalRes, error) {
	statement := `SELECT * FROM localres WHERE patient=$1;`

	rows, err := db.Query(statement, patient)

	if err != nil {
		return LocalRes{}, err
	}

	defer rows.Close()

	var local LocalRes
	for rows.Next() {
		err = rows.Scan(&local.Patient, &local.Sex, &local.Page, &local.Answers, &local.Duration, &local.Aid)
		if err != nil {
			return LocalRes{}, err
		}

	}

	return local, nil
}

func grade(answers []bool, scale Scale, sex string, results *MMPIResults) error {
	if scale.Gender != sex && len(scale.Gender) > 0 {
		return nil
	}

	rawScore := scale.BaseScore

	for _, asw := range scale.Answers {
		for i := 0; i < len(asw); i += 2 {
			if i+1 < len(asw) {
				question := int(asw[i].(float64))
				answer := asw[i+1].(bool)

				if answers[question-1] == answer {
					if len(asw)%2 == 1 {
						rawScore += int32(asw[len(asw)-1].(float64))
					} else {
						rawScore++
					}
				}
			}
		}

	}

	for _, sbs := range scale.SubScales {
		grade(answers, sbs, sex, results)
	}

	var finalScore int32

	if sex == "male" {
		if len(scale.TScores.Male) > 0 {
			finalScore = calculateTScore(rawScore, scale, sex, results)
		} else {
			finalScore = rawScore
		}
	} else {
		if len(scale.TScores.Female) > 0 {
			finalScore = calculateTScore(rawScore, scale, sex, results)
		} else {
			finalScore = rawScore
		}
	}

	var purpose string
	if len(scale.Text) > 0 {
		purpose = scale.Text
	} else if len(scale.Comment) > 0 {
		purpose = scale.Comment
	} else {
		purpose = "N/A"
	}

	results.Categories[len(results.Categories)-1].Scales = append(results.Categories[len(results.Categories)-1].Scales, ScaleResult{ScaleName: scale.Name, ScaleDescription: scale.Title, ScalePupose: purpose, Score: finalScore})

	return nil
}

func calculateTScore(rawScore int32, scale Scale, sex string, results *MMPIResults) int32 {
	var tScores []int32

	if sex == "male" {
		tScores = scale.TScores.Male
	} else {
		tScores = scale.TScores.Female
	}

	if scale.KCorrection > 0.0 {
		var K float32
		var err error
		for _, ct := range results.Categories {
			for _, sr := range ct.Scales {
				if sr.ScaleName == "K" {
					K = float32(sr.Score)
					err = nil
					break
				}
				err = fmt.Errorf("no indications")
			}
		}
		if err != nil {
			K = 1 * scale.KCorrection
		}

		kScore := K + float32(rawScore)

		rawScore = int32(math.Floor(float64(kScore) + 0.5))
	}

	if sex == "male" {
		rawScore -= scale.ScoreOffsets.Male
	} else {
		rawScore -= scale.ScoreOffsets.Female
	}

	return tScores[int(math.Max(0, math.Min(float64(rawScore), float64(len(tScores)-1))))]

}

func (cr *MMPICategoryResult) deriveIndications(scalesData MMPIScales, currentResults *MMPIResults) {
	var indres []string = make([]string, 0)

	for _, sr := range cr.Scales {
		var indications Indications

		for _, category := range scalesData {
			for _, item := range category.Items {
				if item.Name == sr.ScaleName {
					indications = item.Indications
					break
				}

			}
		}

		var F int32 = 0
		var VRIN int32 = 0
		var TRIN int32 = 0
		for _, ct := range currentResults.Categories {
			for _, sv := range ct.Scales {
				if sv.ScaleName == "F" {
					F = sr.Score

				}

				if sv.ScaleName == "VRIN" {
					VRIN = sr.Score

				}

				if sv.ScaleName == "TRIN" {
					TRIN = sr.Score

				}
			}
		}

		switch sr.ScaleName {
		case "VRIN":
			if sr.Score >= 80 {
				indres = append(indres, indications.VRIN__80...)
			}

			if sr.Score < 40 {
				indres = append(indres, indications.VRIN_40...)
			}

		case "TRIN":
			if sr.Score <= -80 || sr.Score >= 80 {
				indres = append(indres, indications.TRIN___80____TRIN__80...)
			}

		case "Hs":
			if sr.Score >= 75 {
				indres = append(indres, indications.Hs__75...)
			}

			if sr.Score >= 55 && sr.Score <= 74 {
				indres = append(indres, indications.Five5__Hs__74...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Six5__Hs__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Five5__Hs__64...)
			}

		case "Fb":
			if sr.Score > F+20 {
				indres = append(indres, indications.Fb_F_20...)
			}

		case "Fp":
			if sr.Score >= 100 && VRIN < 70 && TRIN > 70 {
				indres = append(indres, indications.Fp__100____VRIN_70____TRIN_70)
			}

		case "D":
			if sr.Score >= 75 {
				indres = append(indres, indications.D__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Six5__D__74...)
			}

		case "Hy":
			if sr.Score >= 75 {
				indres = append(indres, indications.Hy__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Six5__Hy__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Five5__Hy__64...)
			}

		case "Pd":
			if sr.Score >= 75 {
				indres = append(indres, indications.Pd__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Six5__Pd__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Five5__Pd__64...)
			}

		case "Mf":
			if sr.Score >= 65 {
				indres = append(indres, indications.Mf__65...)
			}

			if sr.Score >= 45 && sr.Score <= 64 {
				indres = append(indres, indications.Four5__Mf__64...)
			}

			if sr.Score < 45 {
				indres = append(indres, indications.Mf_45...)
			}

		case "Pa":
			if sr.Score >= 75 {
				indres = append(indres, indications.Pa__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Six5__Pa__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Five5__Pa__64...)
			}

		case "Pt":
			if sr.Score >= 75 {
				indres = append(indres, indications.Pt__75...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Five5__Pt__64...)
			}

		case "Sc":
			if sr.Score >= 75 {
				indres = append(indres, indications.Sc__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Six5__Sc__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Five5__Sc__64...)
			}

		case "Ma":
			if sr.Score >= 75 {
				indres = append(indres, indications.Ma__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Ma__65__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Ma__55__64...)
			}

		case "Si":
			if sr.Score >= 75 {
				indres = append(indres, indications.Si__75...)
			}

			if sr.Score >= 65 && sr.Score <= 74 {
				indres = append(indres, indications.Si__65__74...)
			}

			if sr.Score >= 55 && sr.Score <= 64 {
				indres = append(indres, indications.Si__55__64...)
			}

			if sr.Score < 45 {
				indres = append(indres, indications.Si_45...)
			}

		}
	}

	if len(indres) > 0 {
		cr.DerivedIndications = append(cr.DerivedIndications, indres...)
	} else {
		cr.DerivedIndications = append(cr.DerivedIndications, "No particular indications")
	}

}
