// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatModerator

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewLiveChatModerator(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ILiveChatModerator[youtube.LiveChatModerator]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"mod1", "mod2"}),
					WithLiveChatId("liveChat123"),
					WithModeratorChannelId("channel456"),
					WithMaxResults(50),
					WithParts([]string{"snippet", "id"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &LiveChatModerator{
				Fields: common.Fields{
					Service:    svc,
					Parts:      []string{"snippet", "id"},
					Output:     "json",
					MaxResults: 50,
					Ids:        []string{"mod1", "mod2"},
				},
				LiveChatId:         "liveChat123",
				ModeratorChannelId: "channel456",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &LiveChatModerator{Fields: common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &LiveChatModerator{
				Fields: common.Fields{MaxResults: math.MaxInt64},
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &LiveChatModerator{
				Fields: common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithLiveChatId(""),
					WithModeratorChannelId(""),
				},
			},
			want: &LiveChatModerator{
				Fields:             common.Fields{},
				LiveChatId:         "",
				ModeratorChannelId: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithLiveChatId("chat789"),
					WithMaxResults(25),
					WithParts([]string{"id"}),
				},
			},
			want: &LiveChatModerator{
				Fields: common.Fields{
					Parts:      []string{"id"},
					MaxResults: 25,
				},
				LiveChatId: "chat789",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewLiveChatModerator(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewLiveChatModerator() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestLiveChatModerator_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get moderators for live chat",
			opts: []Option{
				WithLiveChatId("live-chat-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("liveChatId") != "live-chat-id" {
					t.Errorf(
						"expected liveChatId=live-chat-id, got %s",
						q.Get("liveChatId"),
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
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write(
								[]byte(`{
					"items": [
						{
							"id": "mod-1",
							"snippet": {
								"liveChatId": "live-chat-id",
								"moderatorDetails": {
									"channelId": "channel-1",
									"displayName": "Moderator 1"
								}
							}
						}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				m := NewLiveChatModerator(opts...)
				got, err := m.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("LiveChatModerator.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("LiveChatModerator.Get() got length = %v, want %v", len(got), tt.wantLen)
				}
			},
		)
	}
}

func TestLiveChatModerator_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("moderator"))

	m := NewLiveChatModerator(
		WithService(svc),
		WithLiveChatId("live-chat-id"),
		WithMaxResults(22),
	)
	got, err := m.Get()
	if err != nil {
		t.Errorf("LiveChatModerator.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("LiveChatModerator.Get() got length = %v, want 22", len(got))
	}
}

func TestLiveChatModerator_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "mod-1",
				"snippet": {
					"liveChatId": "live-chat-id",
					"moderatorDetails": {
						"channelId": "channel-1",
						"displayName": "Moderator 1"
					}
				}
			}
		]
	}`

	common.RunListTest(
		t, mockResponse,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			m := NewLiveChatModerator(
				WithService(svc),
				WithLiveChatId("live-chat-id"),
				WithOutput(output),
				WithMaxResults(1),
			)
			return m.List
		},
	)
}

func TestLiveChatModerator_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert live chat moderator",
			opts: []Option{
				WithLiveChatId("live-chat-id"),
				WithModeratorChannelId("moderator-channel-id"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						LiveChatId       string `json:"liveChatId"`
						ModeratorDetails struct {
							ChannelId string `json:"channelId"`
						} `json:"moderatorDetails"`
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
				if body.Snippet.ModeratorDetails.ChannelId != "moderator-channel-id" {
					t.Errorf(
						"expected snippet.moderatorDetails.channelId=moderator-channel-id, got %s",
						body.Snippet.ModeratorDetails.ChannelId,
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
							_, _ = w.Write([]byte(`{"id": "new-mod-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				m := NewLiveChatModerator(opts...)
				var buf bytes.Buffer
				if err := m.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatModerator.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveChatModerator_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete live chat moderator",
			opts: []Option{
				WithIds([]string{"mod-id"}),
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
				m := NewLiveChatModerator(opts...)
				var buf bytes.Buffer
				if err := m.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatModerator.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}