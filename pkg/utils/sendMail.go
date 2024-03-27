package utils

import (
	"log"
	"net/smtp"
	"os"
)

func SendMail(userEmail string, subject string, body string) error {

	from := os.Getenv("APP_EMAIL_ADDRESS")
	password := os.Getenv("APP_EMAIL_PASSWORD")
	to := userEmail

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	// from := os.Getenv("SMTP_FROM")
	// password := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	message := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
