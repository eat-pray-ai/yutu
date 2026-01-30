// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestNewThumbnail(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want IThumbnail
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithVideoId("video123"),
					WithFile("/path/to/thumbnail.jpg"),
					WithOutput("json"),
					WithJsonpath("id"),
					WithService(svc),
				},
			},
			want: &Thumbnail{
				Fields: &common.Fields{
					Service:  svc,
					Output:   "json",
					Jsonpath: "id",
				},
				VideoId: "video123",
				File:    "/path/to/thumbnail.jpg",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Thumbnail{Fields: &common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithVideoId(""),
					WithFile(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &Thumbnail{
				Fields: &common.Fields{
					Output:   "",
					Jsonpath: "",
				},
				VideoId: "",
				File:    "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithVideoId("myVideo123"),
					WithFile("/images/thumb.png"),
				},
			},
			want: &Thumbnail{
				Fields:  &common.Fields{},
				VideoId: "myVideo123",
				File:    "/images/thumb.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewThumbnail(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewThumbnail() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestThumbnail_Set(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "set thumbnail",
			opts: []Option{
				WithVideoId("video-id"),
				WithFile("test_thumbnail.jpg"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("videoId") != "video-id" {
					t.Errorf(
						"expected videoId=video-id, got %s", r.URL.Query().Get("videoId"),
					)
				}
			},
			wantErr: false,
		},
	}

	err := os.WriteFile("test_thumbnail.jpg", []byte("dummy image content"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer func() {
		_ = os.Remove("test_thumbnail.jpg")
	}()

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write(
								[]byte(`{
					"items": [
						{"default": {"url": "https://example.com/thumb.jpg"}}
					]
				}`),
							)
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				thumb := NewThumbnail(opts...)
				var buf bytes.Buffer
				if err := thumb.Set(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Thumbnail.Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
