package slack_test

import (
	"github.com/int128/slack"
	"log"
)

const webhook = "https://hooks.slack.com/services/..."

func ExampleSend() {
	if err := slack.Send(webhook, &slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	}); err != nil {
		log.Fatalf("Could not send the message to Slack: %s", err)
	}
}
