// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"testing"

	"github.com/spf13/cobra"
)

// YouTube rejects captions.insert (invalidMetadata) for every value of
// snippet.audioTrackType and snippet.trackKind, so insert must not offer them.
func TestInsert_DoesNotOfferTrackFields(t *testing.T) {
	c := newInsertCmd()
	for _, f := range []string{"audioTrackType", "trackKind"} {
		if c.Flags().Lookup(f) != nil {
			t.Errorf("caption insert still has --%s", f)
		}
	}
	for _, p := range []string{"audio_track_type", "track_kind"} {
		if _, ok := insertInSchema.Properties[p]; ok {
			t.Errorf("MCP insert schema still has %s", p)
		}
	}
}

// On update, omitted track flags must stay empty so they don't overwrite the
// caption's current values; passed values must be kept.
func TestUpdate_TrackFieldsResetOnlyWhenOmitted(t *testing.T) {
	c := newUpdateCmd()
	for _, f := range []string{"audioTrackType", "trackKind"} {
		if d := c.Flags().Lookup(f).DefValue; d != "" {
			t.Errorf("update --%s default = %q, want empty", f, d)
		}
	}

	var opts updateOpts
	c2 := &cobra.Command{}
	opts.bindFlags(c2)
	if err := c2.ParseFlags([]string{"--videoId", "v", "--trackKind", "ASR"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	got := opts.build(c2.Flags())
	if got.AudioTrackType != "" {
		t.Errorf("AudioTrackType = %q, want empty when --audioTrackType is omitted", got.AudioTrackType)
	}
	if got.TrackKind != "ASR" {
		t.Errorf("TrackKind = %q, want the passed value ASR", got.TrackKind)
	}
}
