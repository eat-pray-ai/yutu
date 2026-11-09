package channel

import "testing"

// Omitted optional flags must reach Update as nil, or an unrelated update
// (say, only --description) overwrites the channel's keywords with "".
func TestPersistentPreRun_ResetsOmittedKeywords(t *testing.T) {
	keywords = new("")
	if err := updateCmd.ParseFlags([]string{"--id", "channel-id", "--description", "d"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	channelCmd.PersistentPreRun(updateCmd, nil)
	if keywords != nil {
		t.Errorf("keywords = %q, want nil when --keywords is omitted", *keywords)
	}
	if description == nil || *description != "d" {
		t.Errorf("description must keep the passed value")
	}
}
