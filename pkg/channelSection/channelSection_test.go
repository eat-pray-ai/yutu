// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestNewChannelSection(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IChannelSection[youtube.ChannelSection]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"section1", "section2"}),
					WithChannelId("channel123"),
					WithHl("en"),
					WithMine(&mineTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithService(svc),
				},
			},
			want: &ChannelSection{
				Fields:                 &common.Fields{Service: svc},
				Ids:                    []string{"section1", "section2"},
				ChannelId:              "channel123",
				Hl:                     "en",
				Mine:                   &mineTrue,
				OnBehalfOfContentOwner: "owner123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &ChannelSection{Fields: &common.Fields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
				},
			},
			want: &ChannelSection{Fields: &common.Fields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
				},
			},
			want: &ChannelSection{
				Fields: &common.Fields{},
				Mine:   &mineFalse,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithHl(""),
					WithOnBehalfOfContentOwner(""),
				},
			},
			want: &ChannelSection{
				Fields:                 &common.Fields{},
				ChannelId:              "",
				Hl:                     "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"section1"}),
					WithChannelId("partialChannel"),
					WithHl("fr"),
				},
			},
			want: &ChannelSection{
				Fields:    &common.Fields{},
				Ids:       []string{"section1"},
				ChannelId: "partialChannel",
				Hl:        "fr",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannelSection(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewChannelSection() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}

func TestChannelSection_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get channel sections by channelId",
			opts: []Option{
				WithChannelId("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("channelId") != "channel-id" {
					t.Errorf(
						"expected channelId=channel-id, got %s",
						r.URL.Query().Get("channelId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channel sections by id",
			opts: []Option{
				WithIds([]string{"id1", "id2"}),
			},
			verify: func(r *http.Request) {
				ids := r.URL.Query()["id"]
				if len(ids) == 1 && ids[0] == "id1,id2" {
					return
				}
				if len(ids) == 2 && ids[0] == "id1" && ids[1] == "id2" {
					return
				}
				t.Logf("ids: %v", ids)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channel sections mine",
			opts: []Option{
				func(cs *ChannelSection) {
					b := true
					cs.Mine = &b
				},
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("mine") != "true" {
					t.Errorf("expected mine=true, got %s", r.URL.Query().Get("mine"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channel sections by hl",
			opts: []Option{
				WithHl("en"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "en" {
					t.Errorf("expected hl=en, got %s", r.URL.Query().Get("hl"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channel sections with onBehalfOfContentOwner",
			opts: []Option{
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
	}

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
						{"id": "section-1", "snippet": {"title": "Section 1"}}
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
				cs := NewChannelSection(opts...)
				got, err := cs.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"ChannelSection.Get() error = %v, wantErr %v", err, tt.wantErr,
					)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"ChannelSection.Get() got length = %v, want %v", len(got),
						tt.wantLen,
					)
				}
			},
		)
	}
}

func TestChannelSection_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "section-1",
					"snippet": {
						"channelId": "channel-1",
						"title": "Section 1"
					}
				}
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

	tests := []struct {
		name    string
		opts    []Option
		output  string
		wantErr bool
	}{
		{
			name: "list channel sections json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithChannelId("channel-1"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list channel sections yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithChannelId("channel-1"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list channel sections table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithChannelId("channel-1"),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				cs := NewChannelSection(tt.opts...)
				var buf bytes.Buffer
				if err := cs.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"ChannelSection.List() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
				if buf.Len() == 0 {
					t.Errorf("ChannelSection.List() output is empty")
				}
			},
		)
	}
}

func TestChannelSection_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete channel section",
			opts: []Option{
				WithIds([]string{"section-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "delete channel section with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"section-id"}),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.WriteHeader(http.StatusNoContent)
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
				cs := NewChannelSection(opts...)
				var buf bytes.Buffer
				if err := cs.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"ChannelSection.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
