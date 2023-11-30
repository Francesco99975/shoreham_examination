package models

type Exam string

const (
	ASQ Exam = "asq"
	BAI Exam = "bai"
	BDI Exam = "bdi"
	P3  Exam = "p3"
)

const ASQ_MAX_SCORE int16 = 38
const BAI_MAX_SCORE int16 = 57
const BDI_MAX_SCORE int16 = 30
const P3_MAX_SCORE int16 = 85
const P3_ADJUST int16 = 44 // To Be subtracted from p3 score

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
	Categories []MMPICategoryResult
	Duration   int64
}

type MMPIScales []struct {
	Items []struct {
		Answers     [][]interface{} `json:"answers"`
		BaseScore   int64           `json:"baseScore"`
		Code        string          `json:"code"`
		Comment     string          `json:"comment"`
		Gender      string          `json:"gender"`
		Indications struct {
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
		} `json:"indications"`
		KCorrection  float64 `json:"kCorrection"`
		Name         string  `json:"name"`
		ScoreOffsets struct {
			Female int64 `json:"female"`
			Male   int64 `json:"male"`
		} `json:"scoreOffsets"`
		SubScales []struct {
			Answers      [][]interface{} `json:"answers"`
			Name         string          `json:"name"`
			ScoreOffsets struct {
				Female int64 `json:"female"`
				Male   int64 `json:"male"`
			} `json:"scoreOffsets"`
			TScores struct {
				Female []int64 `json:"female"`
				Male   []int64 `json:"male"`
			} `json:"tScores"`
			Title string `json:"title"`
		} `json:"subScales"`
		TScores struct {
			Female []int64 `json:"female"`
			Male   []int64 `json:"male"`
		} `json:"tScores"`
		Text  string `json:"text"`
		Title string `json:"title"`
	} `json:"items"`
	Title string `json:"title"`
}
