// Package slack provides a client for Slack Incoming Webhooks API.
// See https://api.slack.com/docs/messages for details of the API.
package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

type triState *bool

var enable = true
var disable = false

// Enable represents a pointer to true value, for *bool.
var Enable = triState(&enable)

// Disable represents a pointer to false value, for *bool.
var Disable = triState(&disable)

// Message represents a message sent via Incoming Webhooks API.
//
// See https://api.slack.com/docs/message-formatting and
// https://api.slack.com/docs/message-link-unfurling.
type Message struct {
	Username    string       `json:"username,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	Text        string       `json:"text,omitempty"`
	Mrkdwn      triState     `json:"mrkdwn,omitempty"`       // Set false to disable formatting.
	UnfurlMedia triState     `json:"unfurl_media,omitempty"` // Set false to disable unfurling.
	UnfurlLinks triState     `json:"unfurl_links,omitempty"` // Set true to enable unfurling.
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents an attachment of a message.
// See https://api.slack.com/docs/message-attachments for details.
type Attachment struct {
	Fallback   string             `json:"fallback,omitempty"`
	Color      string             `json:"color,omitempty"`
	Pretext    string             `json:"pretext,omitempty"`
	AuthorName string             `json:"author_name,omitempty"`
	AuthorLink string             `json:"author_link,omitempty"`
	AuthorIcon string             `json:"author_icon,omitempty"`
	Title      string             `json:"title,omitempty"`
	TitleLink  string             `json:"title_link,omitempty"`
	Text       string             `json:"text,omitempty"`
	Fields     []AttachmentField  `json:"fields,omitempty"`
	Actions    []AttachmentAction `json:"actions,omitempty"`
	ImageURL   string             `json:"image_url,omitempty"`
	ThumbURL   string             `json:"thumb_url,omitempty"`
	Footer     string             `json:"footer,omitempty"`
	FooterIcon string             `json:"footer_icon,omitempty"`
	Timestamp  int64              `json:"ts,omitempty"`
	MrkdwnIn   []string           `json:"mrkdwn_in,omitempty"` // Valid values are pretext, text, fields
}

// AttachmentField represents a field in an attachment.
// See https://api.slack.com/docs/message-attachments for details.
type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}

// AttachmentAction represents an action in an attachment.
// See https://api.slack.com/docs/message-attachments for details.
type AttachmentAction struct {
	Type  string `json:"type,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
	Style string `json:"style,omitempty"`
}

// Client provides a client for Slack Incoming Webhooks API.
type Client struct {
	WebhookURL string       // Webhook URL (mandatory)
	HTTPClient *http.Client // Default to http.DefaultClient
}

// ErrorResponse represents an error response from Slack API.
// See https://api.slack.com/incoming-webhooks#handling_errors for details.
type ErrorResponse interface {
	StatusCode() int // non-2xx status code
	Body() string    // Response body
}

// GetErrorResponse returns ErrorResponse if Slack API returned an error response.
func GetErrorResponse(err error) ErrorResponse {
	var r ErrorResponse
	if xerrors.As(err, &r) {
		return r
	}
	return nil
}

type slackError struct {
	statusCode int
	body       string
}

func (e *slackError) Error() string {
	return fmt.Sprintf("status=%d, body=%s", e.statusCode, e.body)
}

func (e *slackError) StatusCode() int {
	return e.statusCode
}

func (e *slackError) Body() string {
	return e.body
}

// Send sends the message to Slack.
// It returns an error if a HTTP client returned non-2xx status or network error.
func (c *Client) Send(message *Message) error {
	if message == nil {
		return xerrors.New("message is nil")
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(message); err != nil {
		return xerrors.Errorf("could not encode the message to JSON: %w", err)
	}
	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Post(c.WebhookURL, "application/json", &b)
	if err != nil {
		return xerrors.Errorf("could not send the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		return xerrors.Errorf("error from Slack API: %w", &slackError{
			statusCode: resp.StatusCode,
			body:       string(b),
		})
	}
	return nil
}

// Send sends the message to Slack via Incomming Webhooks API.
// It returns an error if a HTTP client returned non-2xx status or network error.
func Send(WebhookURL string, message *Message) error {
	return (&Client{WebhookURL: WebhookURL}).Send(message)
}
