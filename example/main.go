package main

import (
	"log"
	"os"
	"time"

	"github.com/int128/slack"
)

func main() {
	webhook := os.Getenv("SLACK_WEBHOOK")
	if webhook == "" {
		log.Fatalf("Run with environment variable SLACK_WEBHOOK=https://hooks.slack.com/services/...")
	}
	message := slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
		Attachments: []slack.Attachment{
			{
				Title:      "ALERT",
				Text:       "Hello World!",
				AuthorName: "@author",
				AuthorLink: "https://www.example.com",
				Footer:     "Footer",
				Color:      "danger",
				Timestamp:  time.Now().Unix(),
				Actions: []slack.AttachmentAction{
					{
						Type:  "button",
						Text:  "Detail",
						URL:   "https://www.example.com",
						Style: "danger",
					},
				},
			},
		},
	}
	err := slack.Send(webhook, &message)
	if err != nil {
		log.Fatalf("Could not send the message to Slack: %s", err)
	}
	log.Printf("Sent the message %+v", message)
}
