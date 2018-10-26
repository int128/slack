# slack [![CircleCI](https://circleci.com/gh/int128/slack.svg?style=shield)](https://circleci.com/gh/int128/slack)

A library to send messages to Slack and Mattermost via [Incoming Webhook API](https://api.slack.com/docs/messages), written in Go.

See [GoDoc](https://godoc.org/github.com/int128/slack) for more.


## TL;DR

```go
package main

import (
	"github.com/int128/slack"
	"log"
)

const webhook = "https://hooks.slack.com/services/..."

func main() {
	if err := slack.Send(webhook, &slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	}); err != nil {
		log.Fatalf("Could not send the message to Slack: %s", err)
	}
}
```


## Contributions

This is an open source software licensed under Apache-2.0.
Feel free to open issues and pull requests.
