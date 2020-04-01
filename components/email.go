package components

import (
	"flag"

	"gopkg.in/gomail.v2"
)

// SendMail sends email, with recipient
func SendMail(Config map[string]interface{}, subject string, message string, recipients ...string) error {

	//if test unit running, skip this
	if flag.Lookup("test.v") == nil {
		mailer := gomail.NewMessage()
		mailer.SetHeader("From", Config["email"].(string))
		mailer.SetHeader("To", recipients[0])
		if len(recipients) > 1 {
			mailer.SetHeader("Cc", recipients[1])
		}
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/plain", message)

		dialer := gomail.NewPlainDialer(Config["host"].(string),
			Config["port"].(int),
			Config["email"].(string),
			Config["password"].(string))

		err := dialer.DialAndSend(mailer)
		if err != nil {
			return err
		}
	}

	return nil
}
