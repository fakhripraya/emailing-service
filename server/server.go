package server

import (
	"context"

	protos "github.com/fakhripraya/emailing-service/protos/email"

	"github.com/fakhripraya/emailing-service/data"
	"github.com/fakhripraya/emailing-service/entities"
	"github.com/fakhripraya/emailing-service/protos/email"
	"github.com/hashicorp/go-hclog"
	"gopkg.in/gomail.v2"
)

// Mailer is a gRPC server it implements the methods defined by the MailerServer interface
type Mailer struct {
	email.UnimplementedEmailServer
	logger hclog.Logger
	email  *data.Email
	cred   *entities.EmailCredential
}

// NewMailer creates a new mailer server
func NewMailer(logger hclog.Logger, email *data.Email, cred *entities.EmailCredential) *Mailer {
	newMailer := &Mailer{
		logger: logger,
		email:  email,
		cred:   cred}

	return newMailer
}

// SendEmail is a function to send an email based on the EmailRequest
func (mailer *Mailer) SendEmail(ctx context.Context, rr *protos.EmailRequest) (*protos.EmailResponse, error) {
	var to []string

	for _, temp := range rr.To {
		if mailer.email.IsEmailValid(temp) {
			to = append(to, temp)
		}
	}

	if len(to) == 0 {
		return &protos.EmailResponse{
				ErrorCode:    "404",
				ErrorMessage: "Can't find a valid target Email"},
			nil
	}

	// creating new gomail message
	mail := gomail.NewMessage()
	mail.SetHeader("From", mailer.cred.Username)
	mail.SetHeader("To", to...)
	mail.SetHeader("Cc", rr.Cc...)
	mail.SetHeader("Subject", rr.Subject)
	mail.SetBody("text/html", rr.Body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, mailer.cred.Username, mailer.cred.Password)

	// Send the email
	if err := dialer.DialAndSend(mail); err != nil {
		return &protos.EmailResponse{
				ErrorCode:    "404",
				ErrorMessage: err.Error()},
			nil
	}

	return &protos.EmailResponse{
			ErrorCode:    "200",
			ErrorMessage: "Status Accepted"},
		nil
}
