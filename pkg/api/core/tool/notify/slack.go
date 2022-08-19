package notify

import (
	"github.com/doornoc/config-collector/pkg/api/core/tool/config"
	"github.com/slack-go/slack"
)

func NotifyErrorToSlack(error error) {
	url := config.Conf.Controller.SlackWebhook

	payload := &slack.WebhookMessage{
		Text: "Error(config-collector)",
		Attachments: []slack.Attachment{
			{
				Color: "danger",
				Text:  error.Error(),
			},
		},
	}

	slack.PostWebhook(url, payload)
}
