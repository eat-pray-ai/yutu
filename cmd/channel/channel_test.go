// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"slices"
	"testing"

	"github.com/spf13/cobra"
)

// Omitted optional flags must reach Update as nil, or an unrelated update
// (say, only --description) overwrites the channel's keywords.
func TestUpdate_OmitsUnchangedKeywords(t *testing.T) {
	var opts updateOpts
	c := &cobra.Command{}
	opts.bindFlags(c)
	if err := c.ParseFlags([]string{"--id", "channel-id", "--description", "d"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	got := opts.build(c.Flags())
	if got.Keywords != nil {
		t.Errorf("Keywords = %q, want nil when --keywords is omitted", got.Keywords)
	}
	if got.Description == nil || *got.Description != "d" {
		t.Errorf("Description must keep the passed value")
	}
}

func TestUpdate_KeywordsAreCommaSeparated(t *testing.T) {
	var opts updateOpts
	c := &cobra.Command{}
	opts.bindFlags(c)
	if err := c.ParseFlags([]string{"--id", "channel-id", "--keywords", "AI agents,MCP"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	got := opts.build(c.Flags())
	if want := []string{"AI agents", "MCP"}; !slices.Equal(got.Keywords, want) {
		t.Errorf("Keywords = %q, want %q", got.Keywords, want)
	}
}
