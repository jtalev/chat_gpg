package task_queue

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailHandler struct {
	SenderName       string `json:"sender_name"`
	SenderEmail      string `json:"sender_email"`
	RecipientName    string `json:"recipient_name"`
	RecipientEmail   string `json:"recipient_email"`
	Subject          string `json:"subject"`
	PlainTextContent string `json:"plain_text_content"`
	HtmlContent      string `json:"html_content"`
}

func init() {
	RegisterTaskHandler("send_email", &EmailHandler{})
}

func CreateEmailPayload(
	senderName,
	senderEmail,
	recipientName,
	recipientEmail,
	subject,
	plainTextContent,
	htmlContent string,
) EmailHandler {
	return EmailHandler{
		SenderName:       senderName,
		SenderEmail:      senderEmail,
		RecipientName:    recipientName,
		RecipientEmail:   recipientEmail,
		Subject:          subject,
		PlainTextContent: plainTextContent,
		HtmlContent:      htmlContent,
	}
}

func (e *EmailHandler) ProcessTask(task Task, queue chan Task, db *sql.DB) error {
	err := json.Unmarshal(task.Payload, &e)
	if err != nil {
		log.Printf("error unmarshalling payload: %v", err)
		return err
	}
	log.Println(e)
	err = e.SendEmail()
	if err != nil {
		task.Retries++
		if task.Retries >= task.MaxRetries {
			task.Status = "failed"
			log.Printf("email task failed, max retries exceeded: %v", err)
			err = UpdateTask(task, db)
			if err != nil {
				log.Printf("failed to update db task record, task already failed: %v", err)
				return err
			}
			return err
		}
		err = UpdateTask(task, db)
		if err != nil {
			log.Printf("failed to update db task record, task will not be enqueued: %v", err)
			return err
		}
		queue <- task
		return err
	}
	task.Status = "completed"
	err = UpdateTask(task, db)
	if err != nil {
		return err
	}
	return nil
}

func (e *EmailHandler) SendEmail() error {
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
