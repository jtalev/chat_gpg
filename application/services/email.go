package application

import (
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailSender struct {
	SenderName       string
	SenderEmail      string
	RecipientName    string
	RecipientEmail   string
	Subject          string
	PlainTextContent string
	HtmlContent      string
}

func (e *EmailSender) SendEmail() error {
	from := mail.NewEmail(e.SenderName, e.SenderEmail)
	to := mail.NewEmail(e.RecipientName, e.RecipientEmail)
	message := mail.NewSingleEmail(from, e.Subject, to, e.PlainTextContent, e.HtmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		log.Printf("Successfully sent email: %v", response.StatusCode)
		return nil
	}
}
