package channel

import (
	"slices"
	"testing"
)

// Omitted optional flags must reach Update as nil, or an unrelated update
// (say, only --description) overwrites the channel's keywords.
func TestPersistentPreRun_ResetsOmittedKeywords(t *testing.T) {
	keywords = []string{"stale"}
	updateCmd.Flags().Lookup("keywords").Changed = false
	if err := updateCmd.ParseFlags([]string{"--id", "channel-id", "--description", "d"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	channelCmd.PersistentPreRun(updateCmd, nil)
	if keywords != nil {
		t.Errorf("keywords = %q, want nil when --keywords is omitted", keywords)
	}
	if description == nil || *description != "d" {
		t.Errorf("description must keep the passed value")
	}
}

func TestUpdate_KeywordsAreCommaSeparated(t *testing.T) {
	if err := updateCmd.ParseFlags([]string{"--id", "channel-id", "--keywords", "AI agents,MCP"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	channelCmd.PersistentPreRun(updateCmd, nil)
	if want := []string{"AI agents", "MCP"}; !slices.Equal(keywords, want) {
		t.Errorf("keywords = %q, want %q", keywords, want)
	}
}
