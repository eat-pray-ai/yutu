// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

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

func TestNewWatermark(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IWatermark
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("channel123"),
					WithFile("/path/to/watermark.png"),
					WithInVideoPosition("topRight"),
					WithDurationMs(5000),
					WithOffsetMs(1000),
					WithOffsetType("offsetFromStart"),
					WithOnBehalfOfContentOwner("owner123"),
					WithService(svc),
				},
			},
			want: &Watermark{
				Fields:                 &common.Fields{Service: svc},
				ChannelId:              "channel123",
				File:                   "/path/to/watermark.png",
				InVideoPosition:        "topRight",
				DurationMs:             5000,
				OffsetMs:               1000,
				OffsetType:             "offsetFromStart",
				OnBehalfOfContentOwner: "owner123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Watermark{Fields: &common.Fields{}},
		},
		{
			name: "with zero values",
			args: args{
				opts: []Option{
					WithDurationMs(0),
					WithOffsetMs(0),
				},
			},
			want: &Watermark{
				Fields:     &common.Fields{},
				DurationMs: 0,
				OffsetMs:   0,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithFile(""),
					WithInVideoPosition(""),
					WithOffsetType(""),
					WithOnBehalfOfContentOwner(""),
				},
			},
			want: &Watermark{
				Fields:                 &common.Fields{},
				ChannelId:              "",
				File:                   "",
				InVideoPosition:        "",
				OffsetType:             "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("myChannel"),
					WithFile("/watermarks/logo.png"),
					WithInVideoPosition("bottomLeft"),
					WithDurationMs(10000),
				},
			},
			want: &Watermark{
				Fields:          &common.Fields{},
				ChannelId:       "myChannel",
				File:            "/watermarks/logo.png",
				InVideoPosition: "bottomLeft",
				DurationMs:      10000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewWatermark(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewWatermark() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestWatermark_Set(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "set watermark",
			opts: []Option{
				WithChannelId("channel-id"),
				WithFile("test_watermark.jpg"),
				WithInVideoPosition("topRight"),
				WithDurationMs(5000),
				WithOffsetMs(1000),
				WithOffsetType("offsetFromStart"),
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
			name: "set watermark with onBehalfOfContentOwner",
			opts: []Option{
				WithChannelId("channel-id"),
				WithFile("test_watermark.jpg"),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf("expected onBehalfOfContentOwner=owner-id, got %s", r.URL.Query().Get("onBehalfOfContentOwner"))
				}
			},
			wantErr: false,
		},
	}

	err := os.WriteFile("test_watermark.jpg", []byte("dummy image content"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer os.Remove("test_watermark.jpg")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.WriteHeader(http.StatusNoContent)
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
			w := NewWatermark(opts...)
			var buf bytes.Buffer
			if err := w.Set(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Watermark.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWatermark_Unset(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "unset watermark",
			opts: []Option{
				WithChannelId("channel-id"),
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
			name: "unset watermark with onBehalfOfContentOwner",
			opts: []Option{
				WithChannelId("channel-id"),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf("expected onBehalfOfContentOwner=owner-id, got %s", r.URL.Query().Get("onBehalfOfContentOwner"))
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.WriteHeader(http.StatusNoContent)
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
			w := NewWatermark(opts...)
			var buf bytes.Buffer
			if err := w.Unset(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Watermark.Unset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
