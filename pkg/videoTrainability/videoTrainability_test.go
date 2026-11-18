// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoTrainability

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestVideoTrainability_Get(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		statusCode  int
		response    string
		videoId     string
		wantErr     bool
		errContains string
		wantContain []string
	}{
		{
			name:       "table output format",
			output:     "table",
			statusCode: http.StatusOK,
			response: `{
				"kind": "youtube#videoTrainability",
				"etag": "test-etag",
				"videoId": "vid-123",
				"permitted": ["all"]
			}`,
			videoId: "vid-123",
			wantErr: false,
			wantContain: []string{
				"VIDEO ID",
				"PERMITTED",
				"vid-123",
				"all",
			},
		},
		{
			name:       "json output format",
			output:     "json",
			statusCode: http.StatusOK,
			response: `{
				"kind": "youtube#videoTrainability",
				"etag": "test-etag",
				"videoId": "vid-456",
				"permitted": ["none"]
			}`,
			videoId: "vid-456",
			wantErr: false,
			wantContain: []string{
				`"videoId":"vid-456"`,
				`"permitted":["none"]`,
			},
		},
		{
			name:       "yaml output format",
			output:     "yaml",
			statusCode: http.StatusOK,
			response: `{
				"kind": "youtube#videoTrainability",
				"etag": "test-etag",
				"videoId": "vid-789",
				"permitted": ["partnerA", "partnerB"]
			}`,
			videoId: "vid-789",
			wantErr: false,
			wantContain: []string{
				"videoid: vid-789",
				"partnerA",
				"partnerB",
			},
		},
		{
			name:        "api error case",
			output:      "table",
			statusCode:  http.StatusNotFound,
			response:    `{"error": {"code": 404, "message": "Video not found"}}`,
			videoId:     "nonexistent",
			wantErr:     true,
			errContains: "failed to get video trainability",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			svc, err := youtube.NewService(
				context.Background(),
				option.WithEndpoint(server.URL),
				option.WithHTTPClient(server.Client()),
			)
			if err != nil {
				t.Fatalf("failed to create youtube service: %v", err)
			}

			vt := NewVideoTrainability(
				WithVideoId(tt.videoId),
				WithOutput(tt.output),
				WithService(svc),
			)

			var buf bytes.Buffer
			err = vt.List(&buf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Get() error = %q, want containing %q", err.Error(), tt.errContains)
				}
				return
			}

			out := buf.String()
			for _, want := range tt.wantContain {
				if !strings.Contains(out, want) {
					t.Errorf("Get() output does not contain %q\noutput:\n%s", want, out)
				}
			}
		})
	}
}
