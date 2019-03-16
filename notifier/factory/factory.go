package factory

import (
	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/notifier"
	"github.com/blueskan/gopheart/notifier/email"
	"github.com/blueskan/gopheart/notifier/slack"
	"regexp"
	"strconv"
	"strings"
)

type notifierFactory struct {
}

func NewProviderFactory() NotifierFactory {
	return &notifierFactory{}
}

func (nf *notifierFactory) CreateNotifier(name, notifierType string, config config.NotifierService, threshold int) notifier.Notifier {
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
			threshold,
		)
	case "email":
		regex := regexp.MustCompile(`from:(.*?),title:(.*?),recipients:(.*)`)
		res := regex.FindStringSubmatch(config.Extra)

		from := res[1]
		title := res[2]
		recipients := strings.Split(res[3], "|")

		regex = regexp.MustCompile(`smtp:\/\/(.*?):(.*?)@(.*?):([0-9]*)`)
		res = regex.FindStringSubmatch(config.Url)

		username := res[1]
		password := res[2]
		host := res[3]
		port := res[4]

		intPort, _ := strconv.Atoi(port)

		return email.NewEmail(
			host,
			username,
			password,
			title,
			from,
			config.Message,
			intPort,
			recipients,
			threshold,
		)
	default:
		panic("Cannot parse notifier type `" + notifierType + "`")
	}
}
