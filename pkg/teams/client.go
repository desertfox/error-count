package teams

import (
	"log"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify"
)

var (
	webhookUrl string = os.Getenv("EC_TEAMSWEBHOOK")
)

func SendResults(t, s string) {
	mstClient, _ := goteamsnotify.NewClient()
	card := goteamsnotify.NewMessageCard()
	card.Title = t
	card.Text = s

	if err := mstClient.Send(webhookUrl, card); err != nil {
		log.Printf(
			"failed to send message: %v",
			err,
		)
		os.Exit(1)
	}
}
