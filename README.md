# slack [![GoDoc](https://godoc.org/github.com/int128/slack?status.svg)](https://godoc.org/github.com/int128/slack) [![CircleCI](https://circleci.com/gh/int128/slack.svg?style=shield)](https://circleci.com/gh/int128/slack)

This is a Go package for sending messages via [Slack Incoming Webhooks API](https://api.slack.com/docs/messages) and [Mattermost Incoming Webhooks API](https://developers.mattermost.com/integrate/incoming-webhooks/).
It provides dialects for Slack and Mattermost.

See [GoDoc](https://godoc.org/github.com/int128/slack).


## Examples

```go
package main

import (
	"log"

	"github.com/int128/slack"
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

If you are on Google App Engine:

```go
package main

import (
	"log"
	"net/http"

	"github.com/int128/slack"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

const webhook = "https://hooks.slack.com/services/..."

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	hc := urlfetch.Client(ctx)
	sc := slack.Client{
		WebhookURL: webhook,
		HTTPClient: hc,
	}
	if err := sc.Send(&slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	}); err != nil {
		log.Fatalf("Could not send the message to Slack: %s", err)
		http.Error(w, "Could not send the message to Slack", 500)
	}
}

func main() {
	http.HandleFunc("/", handler)
	appengine.Main()
}
```


## Contributions

This is an open source software licensed under Apache-2.0.
Feel free to open issues and pull requests.
