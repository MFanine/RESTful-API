package routes

import (
    "net/http"
    "series/pkg/controllers"
    "go.mongodb.org/mongo-driver/mongo"
)

// SMTPSender is an implementation of the EmailSender interface
type SMTPSender struct {
    // You can add SMTP configuration fields here if needed
}

// SendEmail sends an email using SMTP
func (s *SMTPSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
    // Implement email sending logic using SMTP here
    // Example:
    // smtp.SendMail(...)
    return nil // Replace this with your actual implementation
}

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(client *mongo.Client) {
    emailSender := &SMTPSender{} // Initialize an instance of SMTPSender
    http.HandleFunc("/register", controllers.Register(client, emailSender))
    // You can add other authentication routes here
}
