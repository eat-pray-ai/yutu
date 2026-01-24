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
					"id": "test-channel-id",
					"snippet": {
						"title": "Test Channel"
					}
				}
			]
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
		input   *listIn
		wantErr bool
		verify  func(t *testing.T, output string)
	}{
		{
			name: "json output",
			input: &listIn{
				CategoryId: "test-category",
				ForHandle:  "test-handle",
				Ids:        []string{"test-channel-id"},
				Parts:      []string{"id", "snippet"},
				Output:     "json",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "test-channel-id")
				assert.Contains(t, output, "Test Channel")
			},
		},
		{
			name: "yaml output",
			input: &listIn{
				ForUsername:            "test-user",
				Hl:                     "en",
				ManagedByMe:            jsonschema.Ptr(true),
				Mine:                   jsonschema.Ptr(true),
				MySubscribers:          jsonschema.Ptr(true),
				OnBehalfOfContentOwner: "test-owner",
				Ids:                    []string{"test-channel-id"},
				Parts:                  []string{"id", "snippet"},
				Output:                 "yaml",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "id: test-channel-id")
				assert.Contains(t, output, "title: Test Channel")
			},
		},
		{
			name: "table output",
			input: &listIn{
				Ids:        []string{"test-channel-id"},
				MaxResults: 5,
				Parts:      []string{"id", "snippet"},
				Output:     "table",
			},
			wantErr: false,
			verify: func(t *testing.T, output string) {
				assert.Contains(t, output, "ID")
				assert.Contains(t, output, "TITLE")
				assert.Contains(t, output, "test-channel-id")
				assert.Contains(t, output, "Test Channel")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err = tt.input.call(&buf, channel.WithService(svc))
			if (err != nil) != tt.wantErr {
				t.Errorf("listIn.call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.verify(t, buf.String())
		})
	}
}
