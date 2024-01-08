package helpers

import (
	"fmt"
	"os"

	"github.com/mailgun/mailgun-go"
)

func SendEmail(test string, description string, patient string, attchment string) (bool, error) {
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_APIKEY"))

	m := mg.NewMessage(fmt.Sprintf("Shoreham Examination <noreply@%s>", os.Getenv("MAILGUN_DOMAIN")), fmt.Sprintf("%s Results for %s", test, patient), description, os.Getenv("RECIPIENT"))

	m.AddAttachment(attchment)
	defer os.Remove(attchment)

	if _, _, err := mg.Send(m); err != nil {
		return false, err
	}

	return true, nil
}
