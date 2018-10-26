package slack_test

import (
	"fmt"
	"github.com/int128/slack"
)

const webhook = "https://hooks.slack.com/services/..."

func ExampleSend() {
	if err := slack.Send(webhook, &slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	}); err != nil {
		panic(fmt.Errorf("Could not send the message to Slack: %s", err))
	}
}
