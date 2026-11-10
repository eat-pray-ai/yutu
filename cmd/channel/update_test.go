package channel

import "testing"

// The handle (custom URL) cannot be changed through the YouTube API, so
// channel update must not offer it on the CLI or in the MCP tool schema.
func TestUpdate_DoesNotOfferCustomUrl(t *testing.T) {
	if f := updateCmd.Flags().Lookup("customUrl"); f != nil {
		t.Errorf("channel update still has a --customUrl flag")
	}
	if _, ok := updateInSchema.Properties["custom_url"]; ok {
		t.Errorf("channel update MCP schema still has custom_url")
	}
}
