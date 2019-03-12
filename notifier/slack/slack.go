package slack

import (
	slackhook "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/blueskan/gopheart/notifier"
	"github.com/blueskan/gopheart/provider"
)

type slack struct {
	webhookUrl string
	username   string
	channel    string
	template   string
}

func NewSlack(webhookUrl, username, channel, template string) *slack {
	return &slack{
		webhookUrl: webhookUrl,
		username:   username,
		channel:    channel,
		template:   template,
	}
}

func (s *slack) Notify(statistics provider.Statistics) error {
	payload := slackhook.Payload{
		Text:     notifier.ComposeMessage(s.template, statistics),
		Username: s.username,
		Channel:  s.channel,
	}

	response := slackhook.Send(s.webhookUrl, "", payload)

	if len(response) <= 0 {
		return nil
	}

	return response[0]
}
