package teams

import (
	"log"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
)

func SendResults(webhookUrl, t, s string) {
	mstClient := goteamsnotify.NewClient()
	mstClient.SkipWebhookURLValidationOnSend(true)

	card := goteamsnotify.NewMessageCard()
	card.Title = t
	card.Text = s

	if err := mstClient.Send(webhookUrl, card); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
	}
}
