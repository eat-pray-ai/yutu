package caption

import (
	"testing"

	"github.com/spf13/pflag"
)

// YouTube rejects captions.insert (invalidMetadata) when snippet.audioTrackType
// ("unknown", "primary") or snippet.trackKind ("standard") is sent, so an insert
// that does not pass these flags must not send them.
func TestInsert_OmittedTrackFieldsAreNotSent(t *testing.T) {
	if err := insertCmd.ParseFlags([]string{"--videoId", "v", "--file", "en.srt", "--language", "en"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	if audioTrackType != "" {
		t.Errorf("audioTrackType = %q, want empty when --audioTrackType is omitted", audioTrackType)
	}
	for _, c := range []string{"insert", "update"} {
		cmd := map[string]interface {
			Flags() *pflag.FlagSet
		}{"insert": insertCmd, "update": updateCmd}[c]
		for _, f := range []string{"audioTrackType", "trackKind"} {
			if d := cmd.Flags().Lookup(f).DefValue; d != "" {
				t.Errorf("%s --%s default = %q, want empty (the flags share one variable)", c, f, d)
			}
		}
	}
	for _, p := range []string{"audio_track_type", "track_kind"} {
		if d := insertInSchema.Properties[p].Default; d != nil {
			t.Errorf("MCP schema must not default %s, got %s", p, d)
		}
	}
	if trackKind != "" {
		t.Errorf("trackKind = %q, want empty when --trackKind is omitted", trackKind)
	}
}
