package helpers

import (
	"fmt"
	"os"

	"github.com/mailgun/mailgun-go"
)

func SendEmail(test string, patient string, attchment string) (bool, error) {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_APIKEY"))

	m := mg.NewMessage("Shoreham Examination <noreply@shorehamexamination>", fmt.Sprintf("%s Results for %s", test, patient), "Attched a file with MMPI-2 results", os.Getenv("RECIPIENT"))

	m.AddAttachment(attchment)
	defer os.Remove(attchment)

	if _, _, err := mg.Send(m); err != nil {
		return false, err
	}

	return true, nil
}
