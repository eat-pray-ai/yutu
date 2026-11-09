package caption

import (
	"testing"

	"github.com/spf13/pflag"
)

// YouTube rejects captions.insert when snippet.audioTrackType is sent
// (invalidMetadata, seen with "unknown" and "primary"), so an insert that
// does not pass --audioTrackType must not send one.
func TestInsert_OmittedAudioTrackTypeIsNotSent(t *testing.T) {
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
		if d := cmd.Flags().Lookup("audioTrackType").DefValue; d != "" {
			t.Errorf("%s --audioTrackType default = %q, want empty (the flags share one variable)", c, d)
		}
	}
	if d := insertInSchema.Properties["audio_track_type"].Default; d != nil {
		t.Errorf("MCP schema must not default audio_track_type, got %s", d)
	}
}
