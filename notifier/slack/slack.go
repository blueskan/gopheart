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
	threshold  int
}

func NewSlack(webhookUrl, username, channel, template string, threshold int) *slack {
	return &slack{
		webhookUrl: webhookUrl,
		username:   username,
		channel:    channel,
		template:   template,
		threshold:  threshold,
	}
}

func (s *slack) GetThreshold() int {
	return s.threshold
}

func (s *slack) GetName() string {
	return "slack"
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
