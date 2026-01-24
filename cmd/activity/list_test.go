// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/activity"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"id": "test-activity-id",
					"snippet": {
						"title": "Test Activity",
						"type": "upload",
						"publishedAt": "2024-01-01T00:00:00Z"
					},
					"contentDetails": {
						"upload": {
							"videoId": "test-video-id"
						}
					}
				}
			],
			"nextPageToken": ""
		}`))
	}))
	defer ts.Close()

	svc, err := youtube.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	tests := []struct {
		name    string
		input   *listIn
		wantErr bool
		verify  func(t *testing.T, output string)
	}{
		{
			name: "json output",
			input: &listIn{
				ChannelId: "test-channel-id",
				Parts:     []string{"id", "snippet", "contentDetails"},
				Output:    "json",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "test-activity-id")
				assert.Contains(t, output, "Test Activity")
				assert.Contains(t, output, "upload")
			},
		},
		{
			name: "table output",
			input: &listIn{
				Home:       jsonschema.Ptr(true),
				MaxResults: 10,
				RegionCode: "US",
				Output:     "table",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "ID")
				assert.Contains(t, output, "TITLE")
				assert.Contains(t, output, "TYPE")
				assert.Contains(t, output, "TIME")
				assert.Contains(t, output, "test-activity-id")
				assert.Contains(t, output, "Test Activity")
			},
		},
		{
			name: "yaml output",
			input: &listIn{
				Mine:            jsonschema.Ptr(true),
				PublishedAfter:  "2024-01-01T00:00:00Z",
				PublishedBefore: "2024-12-31T23:59:59Z",
				Output:          "yaml",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "id: test-activity-id")
				assert.Contains(t, output, "title: Test Activity")
				assert.Contains(t, output, "type: upload")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err = tt.input.call(&buf, activity.WithService(svc))
			if (err != nil) != tt.wantErr {
				t.Errorf("listIn.call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(t, buf.String())
		})
	}
}
