package utils

import (
	"log"
	"net/smtp"
)

func SendEmail(to string, verificationToken string) error {
	from := "faninemouad@gmail.com"
	password := "uobefmawbkbjpdps"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	verificationURL := "http://localhost:8080/verify-email?token=" + verificationToken

	message := []byte("To: " + to + "\r\n" +
		"Subject: Email Verification\r\n" +
		"\r\n" +
		"Please click the following link to verify your email: " + verificationURL + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}
