package layers

import (
	"fmt"
	"log"
	"os"
)

type EmailSenderInterface interface {
	Send(recipients []string, sender string, emailBody string) error
}

type emailSender struct {
	logger log.Logger
	mode   string
}

func NewEmailSender(logger log.Logger) EmailSenderInterface {
	mode := os.Getenv("RUN_MODE")
	return &emailSender{
		logger: logger,
		mode:   mode,
	}
}

func (em *emailSender) Send(recipients []string, sender string, emailBody string) error {
	if em.mode == "development" {
		fmt.Printf("The email is sent: %s\n", emailBody)
	}
	em.logger.Printf("Email is sent from %s to %v\n", sender, recipients)
	return nil
}
