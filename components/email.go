package components

import (
	"flag"

	"gopkg.in/gomail.v2"
)

// SendMail sends email
func SendMail(Config map[string]interface{}, to, subject, message string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", Config["email"].(string))
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", message)

	dialer := gomail.NewPlainDialer(Config["host"].(string),
		Config["port"].(int),
		Config["email"].(string),
		Config["password"].(string))

	if flag.Lookup("test.v") == nil {
		err := dialer.DialAndSend(mailer)
		if err != nil {
			return err
		}
	}

	return nil
}
