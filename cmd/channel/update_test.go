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

	var buf bytes.Buffer
	input := &updateIn{
		Ids:    []string{"test-channel-id"},
		Title:  "Updated Channel",
		Output: "json",
	}

	err = input.call(&buf, channel.WithService(svc))

	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "test-channel-id")
	assert.Contains(t, output, "Updated Channel")
}
