package slack_test

import (
	"encoding/json"
	"fmt"
	"testing"

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

func ExampleClient_Send() {
	c := &slack.Client{
		WebhookURL: webhook,
		HTTPClient: nil, // urlfetch.Client(ctx) on App Engine
	}
	if err := c.Send(&slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	}); err != nil {
		panic(fmt.Errorf("Could not send the message to Slack: %s", err))
	}
}

func TestMarshalMessage_empty(t *testing.T) {
	m := slack.Message{}
	b, err := json.Marshal(&m)
	if err != nil {
		t.Fatalf("Could not marshal slack.Message: %s", err)
	}
	s := string(b)
	if want := "{}"; want != s {
		t.Errorf("JSON wants %s but %s", want, s)
	}
}

func TestMarshalMessage_boolFlag(t *testing.T) {
	m := slack.Message{
		Mrkdwn: slack.Disable,
	}
	b, err := json.Marshal(&m)
	if err != nil {
		t.Fatalf("Could not marshal slack.Message: %s", err)
	}
	s := string(b)
	if want := `{"mrkdwn":false}`; want != s {
		t.Errorf("JSON wants %s but %s", want, s)
	}
}
