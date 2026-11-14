package caption

import "testing"

// YouTube rejects captions.insert (invalidMetadata) for every value of
// snippet.audioTrackType and snippet.trackKind, so insert must not offer them.
func TestInsert_DoesNotOfferTrackFields(t *testing.T) {
	for _, f := range []string{"audioTrackType", "trackKind"} {
		if insertCmd.Flags().Lookup(f) != nil {
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
	for _, f := range []string{"audioTrackType", "trackKind"} {
		if d := updateCmd.Flags().Lookup(f).DefValue; d != "" {
			t.Errorf("update --%s default = %q, want empty", f, d)
		}
	}

	audioTrackType, trackKind = "primary", "forced"
	updateCmd.Flags().Lookup("audioTrackType").Changed = false
	updateCmd.Flags().Lookup("trackKind").Changed = false
	if err := updateCmd.ParseFlags([]string{"--videoId", "v", "--trackKind", "ASR"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	captionCmd.PersistentPreRun(updateCmd, nil)
	if audioTrackType != "" {
		t.Errorf("audioTrackType = %q, want empty when --audioTrackType is omitted", audioTrackType)
	}
	if trackKind != "ASR" {
		t.Errorf("trackKind = %q, want the passed value ASR", trackKind)
	}
}
