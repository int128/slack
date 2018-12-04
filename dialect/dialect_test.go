package dialect

import "testing"

func TestSlack_Mention(t *testing.T) {
	var d Slack
	if m := d.Mention(""); m != "" {
		t.Errorf("Mention() wants empty but %s", m)
	}
	if m := d.Mention("foo"); m != "<@foo>" {
		t.Errorf("Mention(foo) wants <@foo> but %s", m)
	}
}

func TestMattermost_Mention(t *testing.T) {
	var d Mattermost
	if m := d.Mention(""); m != "" {
		t.Errorf("Mention() wants empty but %s", m)
	}
	if m := d.Mention("foo"); m != "@foo" {
		t.Errorf("Mention(foo) wants @foo but %s", m)
	}
}
