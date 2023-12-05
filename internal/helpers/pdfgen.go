package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Francesco99975/shorehamex/internal/models"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func GeneratePDF(results models.MMPIResults) (string, error) {
	date := time.Now().Format(time.RubyDate)
	cfg := config.NewBuilder().
		WithPageNumber("Page {current} of {total}", props.RightBottom).
		WithMargins(10, 15, 10).
		Build()

	darkGrayColor := getDarkGrayColor()
	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	err := m.RegisterHeader(getPageHeader(results.Patient, results.Sex, date))
	if err != nil {
		return "", err
	}

	m.AddRows(text.NewRow(25, fmt.Sprintf("MMPI-2 results for %s", results.Patient), props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Size:  20,
	}))

	for _, category := range results.Categories {

		m.AddRows(text.NewRow(20, category.Title, props.Text{
			Top:   1.5,
			Style: fontstyle.Bold,
			Align: align.Center,
			Size:  16,
			Color: &props.WhiteColor,
		}).WithStyle(&props.Cell{BackgroundColor: darkGrayColor}))

		m.AddRows(getTransactions(category.Scales)...)

		m.AddRows(text.NewRow(75, fmt.Sprintf("Indications: %s", strings.Join(category.DerivedIndications, ",")), props.Text{
			Top:   5,
			Style: fontstyle.Italic,
			Align: align.Left,
			Size:  16,
		}))
	}

	document, err := m.Generate()
	if err != nil {
		return "", err
	}

	filename := strings.ReplaceAll(fmt.Sprintf("%s-%s.pdf", results.Patient, date), " ", "_")

	document.Save(filename)

	return filename, err

}

func getTransactions(scaleResults []models.ScaleResult) []core.Row {
	rows := []core.Row{
		row.New(15).Add(
			text.NewCol(2, "Scale", props.Text{Size: 12, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(4, "Scale Description", props.Text{Size: 12, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(4, "Scale Purpose", props.Text{Size: 12, Align: align.Center, Style: fontstyle.Bold}),
			text.NewCol(2, "Scale Score", props.Text{Size: 12, Align: align.Center, Style: fontstyle.Bold}),
		),
	}
	var contentsRow []core.Row

	for i, scale := range scaleResults {
		r := row.New(15).Add(
			text.NewCol(2, scale.ScaleName, props.Text{Size: 12, Align: align.Center}),
			text.NewCol(4, scale.ScaleDescription, props.Text{Size: 12, Align: align.Center}),
			text.NewCol(4, scale.ScalePupose, props.Text{Size: 12, Align: align.Center}),
			text.NewCol(2, strconv.Itoa(int(scale.Score)), props.Text{Size: 12, Align: align.Center}),
		)
		if i%2 == 0 {
			gray := getGrayColor()
			r.WithStyle(&props.Cell{BackgroundColor: gray})
		}

		contentsRow = append(contentsRow, r)
	}

	rows = append(rows, contentsRow...)

	return rows
}

func getPageHeader(patient string, sex string, date string) core.Row {
	return row.New(30).Add(
		col.New(10).Add(
			text.New(fmt.Sprintf("Patient: %s", patient), props.Text{
				Size:  12,
				Align: align.Left,
				Style: fontstyle.BoldItalic,
				Color: getBlueColor(),
			}),
			text.New(fmt.Sprintf("Sex: %s", sex), props.Text{
				Top:   6,
				Size:  12,
				Align: align.Left,
				Color: getBlueColor(),
			}),
			text.New(fmt.Sprintf("Date: %s", date), props.Text{
				Top:   12,
				Size:  12,
				Align: align.Left,
				Color: getBlueColor(),
			}),
		),
	)
}

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}
