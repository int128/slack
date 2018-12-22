package slack_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/int128/slack"
)

const webhook = "https://hooks.slack.com/services/..."

func ExampleSend() {
	err := slack.Send(webhook, &slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	})
	if err != nil {
		panic(fmt.Errorf("could not send the message to Slack: %s", err))
	}
}

func ExampleClient_Send() {
	c := &slack.Client{
		WebhookURL: webhook,
		HTTPClient: nil, // urlfetch.Client(ctx) on App Engine
	}
	err := c.Send(&slack.Message{
		Username:  "mybot",
		IconEmoji: ":star:",
		Text:      "Hello World!",
	})
	if err != nil {
		panic(fmt.Errorf("could not send the message to Slack: %s", err))
	}
}

func TestSend(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != "POST" {
			t.Errorf("Method wants POST but %s", r.Method)
		}
		if r.URL.Path != "/webhook" {
			t.Errorf("Path wants /webhook but %s", r.URL.Path)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Could not read body: %s", err)
			return
		}
		body := strings.TrimSpace(string(b))
		if want := `{"text":"Hello World!"}`; body != want {
			t.Errorf("Body wants %s but %s", want, body)
		}
	}))
	defer s.Close()
	err := slack.Send(s.URL+"/webhook", &slack.Message{Text: "Hello World!"})
	if err != nil {
		t.Fatalf("Send returned error: %s", err)
	}
}

func TestMessage(t *testing.T) {
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

func Test_triState(t *testing.T) {
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
