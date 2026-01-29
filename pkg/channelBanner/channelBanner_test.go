// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestChannelBanner_Insert_Error(t *testing.T) {
	// Setup pkg.Root for file tests
	tmpDir := t.TempDir()
	f, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	oldRoot := pkg.Root
	pkg.Root = f
	defer func() { pkg.Root = oldRoot }()
	defer f.Close()

	svc, _ := youtube.NewService(context.Background(), option.WithAPIKey("test"))

	// Test: File open error
	cb := NewChannelBanner(WithFile("non_existent.jpg"), WithService(svc))
	if err := cb.Insert(&bytes.Buffer{}); err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// Test: API error
	if err := os.WriteFile(tmpDir+"/test.jpg", []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	svc, _ = youtube.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test"))
	cb = NewChannelBanner(WithFile("test.jpg"), WithService(svc), WithChannelId("cid"))
	if err := cb.Insert(&bytes.Buffer{}); err == nil {
		t.Error("expected error from API, got nil")
	}
}

func TestChannelBanner_Insert_Output(t *testing.T) {
	tmpDir := t.TempDir()
	f, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	oldRoot := pkg.Root
	pkg.Root = f
	defer func() { pkg.Root = oldRoot }()
	defer f.Close()

	if err := os.WriteFile(tmpDir+"/test.jpg", []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"url": "http://example.com/banner.jpg"}`))
	}))
	defer ts.Close()
	svc, _ := youtube.NewService(context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test"))

	tests := []struct {
		name        string
		output      string
		expectEmpty bool
	}{
		{"json", "json", false},
		{"yaml", "yaml", false},
		{"silent", "silent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := NewChannelBanner(WithFile("test.jpg"), WithService(svc), WithChannelId("cid"), WithOutput(tt.output))
			var buf bytes.Buffer
			if err := cb.Insert(&buf); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectEmpty {
				if buf.Len() > 0 {
					t.Errorf("expected empty output, got %s", buf.String())
				}
			} else {
				if buf.Len() == 0 {
					t.Error("expected output, got empty")
				}
			}
		})
	}
}

func TestNewChannelBanner(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want IChannelBanner
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("channel123"),
					WithFile("/path/to/banner.jpg"),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &ChannelBanner{
				Fields: &common.Fields{
					Output:   "json",
					Jsonpath: "items.id",
					Service:  svc,
				},
				ChannelId:                     "channel123",
				File:                          "/path/to/banner.jpg",
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &ChannelBanner{Fields: &common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithFile(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &ChannelBanner{
				Fields: &common.Fields{
					Output:   "",
					Jsonpath: "",
				},
				ChannelId:                     "",
				File:                          "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("partialChannel"),
					WithFile("/partial/banner.png"),
					WithOutput("yaml"),
				},
			},
			want: &ChannelBanner{
				Fields:    &common.Fields{Output: "yaml"},
				ChannelId: "partialChannel",
				File:      "/partial/banner.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannelBanner(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewChannelBanner() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestChannelBanner_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert channel banner",
			opts: []Option{
				WithChannelId("channel-id"),
				WithFile("test_banner.jpg"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("channelId") != "channel-id" {
					t.Errorf("expected channelId=channel-id, got %s", r.URL.Query().Get("channelId"))
				}
			},
			wantErr: false,
		},
		{
			name: "insert channel banner with onBehalfOfContentOwner",
			opts: []Option{
				WithChannelId("channel-id"),
				WithFile("test_banner.jpg"),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf("expected onBehalfOfContentOwner=owner-id, got %s", r.URL.Query().Get("onBehalfOfContentOwner"))
				}
			},
			wantErr: false,
		},
		{
			name: "insert channel banner with onBehalfOfContentOwnerChannel",
			opts: []Option{
				WithChannelId("channel-id"),
				WithFile("test_banner.jpg"),
				WithOnBehalfOfContentOwnerChannel("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwnerChannel") != "channel-id" {
					t.Errorf("expected onBehalfOfContentOwnerChannel=channel-id, got %s", r.URL.Query().Get("onBehalfOfContentOwnerChannel"))
				}
			},
			wantErr: false,
		},
	}

	err := os.WriteFile("test_banner.jpg", []byte("dummy image content"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer os.Remove("test_banner.jpg")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"url": "http://example.com/banner.jpg"}`))
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

			opts := append([]Option{WithService(svc)}, tt.opts...)
			cb := NewChannelBanner(opts...)
			var buf bytes.Buffer
			if err := cb.Insert(&buf); (err != nil) != tt.wantErr {
				t.Errorf("ChannelBanner.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
