package email

import (
	"crypto/tls"
	"github.com/blueskan/gopheart/notifier"
	"github.com/blueskan/gopheart/provider"
	"gopkg.in/gomail.v2"
)

type email struct {
	host string
	port int
	username string
	password   string
	title    string
	from   string
	recipients   []string
	message string
	threshold int
}

func NewEmail(host, username, password, title, from, message string, port int, recipients []string, threshold int) *email {
	return &email{
		host: host,
		port: port,
		username: username,
		password: password,
		title: title,
		from: from,
		recipients: recipients,
		message: message,
		threshold: threshold,
	}
}

func (e *email) GetThreshold() int {
	return e.threshold
}

func (e *email) GetName() string {
	return "email"
}

func (e *email) Notify(statistics provider.Statistics) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", e.recipients...)
	m.SetHeader("Subject", notifier.ComposeMessage(e.title, statistics))
	m.SetBody("text/html", notifier.ComposeMessage(e.message, statistics))

	d := gomail.NewPlainDialer(e.host, e.port, e.username, e.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
