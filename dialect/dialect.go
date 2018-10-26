// Package dialect provides functionality for easily switching Slack or Mattermost.
package dialect

import "fmt"

// Dialect provides some formatters.
type Dialect interface {
	Mention(string) string
}

// Slack is a dialect for Slack API.
type Slack struct{}

// Mention returns the user mention for Slack.
func (d *Slack) Mention(username string) string {
	return fmt.Sprintf("<@%s>", username)
}

// Mattermost is a dialect for Mattermost API.
type Mattermost struct{}

// Mention returns the user mention for Mattermost.
func (d *Mattermost) Mention(username string) string {
	return fmt.Sprintf("@%s", username)
}
