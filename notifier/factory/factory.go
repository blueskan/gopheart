package factory

import (
	"strconv"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/notifier"
	"github.com/blueskan/gopheart/notifier/email"
	"github.com/blueskan/gopheart/notifier/slack"
)

type notifierFactory struct {
}

func NewProviderFactory() NotifierFactory {
	return &notifierFactory{}
}

func (nf *notifierFactory) CreateNotifier(name, notifierType string, config config.NotifierService, threshold int) notifier.Notifier {
	switch notifierType {
	case "slack":
		return slack.NewSlack(
			config.SlackURL,
			config.SlackUsername,
			config.SlackChannel,
			config.Message,
			threshold,
		)
	case "email":
		intPort, _ := strconv.Atoi(config.SMTPPort)

		return email.NewEmail(
			config.SMTPHost,
			config.SMTPUsername,
			config.SMTPPassword,
			config.MailTitle,
			config.MailFrom,
			config.Message,
			intPort,
			config.MailRecipients,
			threshold,
		)
	default:
		panic("Cannot parse notifier type `" + notifierType + "`")
	}
}
