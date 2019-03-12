package factory

import (
	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/notifier"
	"github.com/blueskan/gopheart/notifier/slack"
	"regexp"
)

type notifierFactory struct {
}

func NewProviderFactory() NotifierFactory {
	return &notifierFactory{}
}

func (nf *notifierFactory) CreateNotifier(name, notifierType string, config config.NotifierService) notifier.Notifier {
	switch notifierType {
	case "slack":
		regex := regexp.MustCompile(`channel:(.*?),username:(.*?)`)
		res := regex.FindStringSubmatch(config.Extra)

		channel := res[1]
		username := res[2]

		return slack.NewSlack(
			config.Url,
			username,
			channel,
			config.Message,
		)
	default:
		panic("Cannot parse notifier type `" + notifierType + "`")
	}
}
