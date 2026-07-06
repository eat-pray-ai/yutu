// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatBan

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewLiveChatBan(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ILiveChatBan
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"ban1", "ban2"}),
					WithLiveChatId("liveChat123"),
					WithBannedUserChannelId("channel456"),
					WithBanDurationSeconds(300),
					WithBanType("temporary"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &LiveChatBan{
				Fields: common.Fields{
					Service: svc,
					Parts:   []string{"snippet"},
					Output:  "json",
					Ids:     []string{"ban1", "ban2"},
				},
				LiveChatId:          "liveChat123",
				BannedUserChannelId: "channel456",
				BanDurationSeconds:  300,
				BanType:             "temporary",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &LiveChatBan{Fields: common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithLiveChatId(""),
					WithBannedUserChannelId(""),
					WithBanType(""),
				},
			},
			want: &LiveChatBan{
				Fields:              common.Fields{},
				LiveChatId:          "",
				BannedUserChannelId: "",
				BanType:             "",
			},
		},
		{
			name: "with permanent ban",
			args: args{
				opts: []Option{
					WithLiveChatId("chat789"),
					WithBannedUserChannelId("user123"),
					WithBanType("permanent"),
				},
			},
			want: &LiveChatBan{
				Fields:              common.Fields{},
				LiveChatId:          "chat789",
				BannedUserChannelId: "user123",
				BanType:             "permanent",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithLiveChatId("chat456"),
					WithBanType("temporary"),
					WithBanDurationSeconds(600),
				},
			},
			want: &LiveChatBan{
				Fields:             common.Fields{},
				LiveChatId:         "chat456",
				BanType:            "temporary",
				BanDurationSeconds: 600,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewLiveChatBan(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewLiveChatBan() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestLiveChatBan_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert live chat ban",
			opts: []Option{
				WithLiveChatId("live-chat-id"),
				WithBannedUserChannelId("banned-channel-id"),
				WithBanType("permanent"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						LiveChatId        string `json:"liveChatId"`
						BannedUserDetails struct {
							ChannelId string `json:"channelId"`
						} `json:"bannedUserDetails"`
						Type string `json:"type"`
					} `json:"snippet"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Snippet.LiveChatId != "live-chat-id" {
					t.Errorf(
						"expected snippet.liveChatId=live-chat-id, got %s",
						body.Snippet.LiveChatId,
					)
				}
				if body.Snippet.BannedUserDetails.ChannelId != "banned-channel-id" {
					t.Errorf(
						"expected snippet.bannedUserDetails.channelId=banned-channel-id, got %s",
						body.Snippet.BannedUserDetails.ChannelId,
					)
				}
				if body.Snippet.Type != "permanent" {
					t.Errorf("expected snippet.type=permanent, got %s", body.Snippet.Type)
				}
			},
			wantErr: false,
		},
		{
			name: "insert temporary ban with duration",
			opts: []Option{
				WithLiveChatId("live-chat-id"),
				WithBannedUserChannelId("banned-channel-id"),
				WithBanType("temporary"),
				WithBanDurationSeconds(300),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						BanDurationSeconds string `json:"banDurationSeconds"`
						Type               string `json:"type"`
					} `json:"snippet"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Snippet.Type != "temporary" {
					t.Errorf("expected snippet.type=temporary, got %s", body.Snippet.Type)
				}
				if body.Snippet.BanDurationSeconds != "300" {
					t.Errorf(
						"expected snippet.banDurationSeconds=300, got %s",
						body.Snippet.BanDurationSeconds,
					)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"id": "new-ban-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveChatBan(opts...)
				var buf bytes.Buffer
				if err := b.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatBan.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveChatBan_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete live chat ban",
			opts: []Option{
				WithIds([]string{"ban-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.WriteHeader(http.StatusNoContent)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveChatBan(opts...)
				var buf bytes.Buffer
				if err := b.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatBan.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
