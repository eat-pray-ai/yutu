// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/channel"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			w.Write([]byte(`{
				"items": [
					{
						"id": "test-channel-id",
						"snippet": {
							"title": "Original Channel"
						}
					}
				]
			}`))
			return
		}

		// Response for Update call
		w.Write([]byte(`{
			"id": "test-channel-id",
			"snippet": {
				"title": "Updated Channel"
			}
		}`))
	}))
	defer ts.Close()

	svc, err := youtube.NewService(
		context.Background(),
		option.WithEndpoint(ts.URL),
		option.WithAPIKey("test-key"),
	)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	tests := []struct {
		name    string
		input   *updateIn
		wantErr bool
		verify  func(t *testing.T, output string)
	}{
		{
			name: "json output",
			input: &updateIn{
				Ids:             []string{"test-channel-id"},
				Country:         "US",
				CustomUrl:       "my-channel",
				DefaultLanguage: "en",
				Title:           "Updated Channel",
				Output:          "json",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "test-channel-id")
				assert.Contains(t, output, "Updated Channel")
			},
		},
		{
			name: "yaml output",
			input: &updateIn{
				Ids:         []string{"test-channel-id"},
				Description: "New Description",
				Title:       "Updated Channel",
				Output:      "yaml",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "id: test-channel-id")
				assert.Contains(t, output, "title: Updated Channel")
			},
		},
		{
			name: "silent output",
			input: &updateIn{
				Ids:    []string{"test-channel-id"},
				Title:  "Updated Channel",
				Output: "silent",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Empty(t, output)
			},
		},
		{
			name: "default output",
			input: &updateIn{
				Ids:    []string{"test-channel-id"},
				Title:  "Updated Channel",
				Output: "",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "Channel updated: test-channel-id")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err = tt.input.call(&buf, channel.WithService(svc))
			if (err != nil) != tt.wantErr {
				t.Errorf("updateIn.call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(t, buf.String())
		})
	}
}
