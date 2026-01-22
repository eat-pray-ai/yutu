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

	var buf bytes.Buffer
	input := &listIn{
		Home:   jsonschema.Ptr(true),
		Mine:   jsonschema.Ptr(true),
		Parts:  []string{"id", "snippet", "contentDetails"},
		Output: "json",
	}

	err = input.call(&buf, activity.WithService(svc))

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "test-activity-id")
	assert.Contains(t, output, "Test Activity")
	assert.Contains(t, output, "upload")
}
