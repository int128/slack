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

func TestSend_NilMessage(t *testing.T) {
	err := slack.Send("/webhook", nil)
	if err == nil {
		t.Error("err wants non-nil but nil")
	}
}

func TestSend_EmptyWebhookURL(t *testing.T) {
	err := slack.Send("", &slack.Message{})
	if err == nil {
		t.Error("err wants non-nil but nil")
	}
}

func ExampleGetErrorResponse() {
	err := slack.Send(webhook, &slack.Message{
		Text: "Hello World!",
	})
	if err != nil {
		if resp := slack.GetErrorResponse(err); resp != nil {
			if resp.StatusCode() >= 500 {
				// you can retry sending the message
			}
		}
		panic(fmt.Errorf("could not send the message to Slack: %s", err))
	}
}

func TestGetErrorResponse(t *testing.T) {
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
		if body != "{}" {
			t.Errorf("Body wants {} but %s", body)
		}
		w.WriteHeader(400)
		if _, err := fmt.Fprint(w, "invalid_payload"); err != nil {
			t.Errorf("Could not write body: %s", err)
		}
	}))
	defer s.Close()
	err := slack.Send(s.URL+"/webhook", &slack.Message{})
	if err == nil {
		t.Fatalf("err wants non-nil but got nil")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("err.Error should contain status code but %s", err.Error())
	}
	if !strings.Contains(err.Error(), "invalid_payload") {
		t.Errorf("err.Error should contain body but %s", err.Error())
	}
	errResp := slack.GetErrorResponse(err)
	if errResp == nil {
		t.Fatalf("GetErrorResponse wants non-nil but nil")
	}
	if errResp.StatusCode() != 400 {
		t.Errorf("StatusCode wants 400 but %d", errResp.StatusCode())
	}
	if errResp.Body() != "invalid_payload" {
		t.Errorf("Body wants invalid_payload but %s", errResp.Body())
	}
}

func TestGetErrorResponse_NonErrorResponse(t *testing.T) {
	err := fmt.Errorf("some error")
	resp := slack.GetErrorResponse(err)
	if resp != nil {
		t.Errorf("GetErrorResponse wants nil but got %+v", resp)
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
