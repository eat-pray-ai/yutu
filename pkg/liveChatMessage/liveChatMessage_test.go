// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveChatMessage

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

func TestNewLiveChatMessage(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ILiveChatMessage[youtube.LiveChatMessage]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"msg1", "msg2"}),
					WithLiveChatId("liveChat123"),
					WithMessageText("hello world"),
					WithStatus("closed"),
					WithMaxResults(50),
					WithParts([]string{"snippet", "id", "authorDetails"}),
					WithOutput("json"),
					WithService(svc),
					WithHl("en"),
				},
			},
			want: &LiveChatMessage{
				Fields: common.Fields{
					Service:    svc,
					Parts:      []string{"snippet", "id", "authorDetails"},
					Output:     "json",
					MaxResults: 50,
					Ids:        []string{"msg1", "msg2"},
					Hl:         "en",
				},
				LiveChatId:  "liveChat123",
				MessageText: "hello world",
				Status:      "closed",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &LiveChatMessage{Fields: common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &LiveChatMessage{
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
			want: &LiveChatMessage{
				Fields: common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithLiveChatId(""),
					WithMessageText(""),
					WithStatus(""),
				},
			},
			want: &LiveChatMessage{
				Fields:      common.Fields{},
				LiveChatId:  "",
				MessageText: "",
				Status:      "",
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
			want: &LiveChatMessage{
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
				if got := NewLiveChatMessage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewLiveChatMessage() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}

func TestLiveChatMessage_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get messages for live chat",
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
							"id": "msg-1",
							"snippet": {
								"liveChatId": "live-chat-id",
								"type": "textMessageEvent",
								"displayMessage": "Hello!"
							},
							"authorDetails": {
								"channelId": "channel-1",
								"displayName": "User 1"
							}
						}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				m := NewLiveChatMessage(opts...)
				got, err := m.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatMessage.Get() error = %v, wantErr %v", err, tt.wantErr,
					)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"LiveChatMessage.Get() got length = %v, want %v", len(got),
						tt.wantLen,
					)
				}
			},
		)
	}
}

func TestLiveChatMessage_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("message"))

	m := NewLiveChatMessage(
		WithService(svc),
		WithLiveChatId("live-chat-id"),
		WithMaxResults(22),
	)
	got, err := m.Get()
	if err != nil {
		t.Errorf("LiveChatMessage.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("LiveChatMessage.Get() got length = %v, want 22", len(got))
	}
}

func TestLiveChatMessage_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "msg-1",
				"snippet": {
					"liveChatId": "live-chat-id",
					"type": "textMessageEvent",
					"displayMessage": "Hello!"
				},
				"authorDetails": {
					"channelId": "channel-1",
					"displayName": "User 1"
				}
			}
		]
	}`

	common.RunListTest(
		t, mockResponse,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			m := NewLiveChatMessage(
				WithService(svc),
				WithLiveChatId("live-chat-id"),
				WithOutput(output),
				WithMaxResults(1),
			)
			return m.List
		},
	)
}

func TestLiveChatMessage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert live chat message",
			opts: []Option{
				WithLiveChatId("live-chat-id"),
				WithMessageText("Hello, world!"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						LiveChatId         string `json:"liveChatId"`
						TextMessageDetails struct {
							MessageText string `json:"messageText"`
						} `json:"textMessageDetails"`
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
				if body.Snippet.TextMessageDetails.MessageText != "Hello, world!" {
					t.Errorf(
						"expected snippet.textMessageDetails.messageText=Hello, world!, got %s",
						body.Snippet.TextMessageDetails.MessageText,
					)
				}
				if body.Snippet.Type != "textMessageEvent" {
					t.Errorf(
						"expected snippet.type=textMessageEvent, got %s", body.Snippet.Type,
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
							_, _ = w.Write([]byte(`{"id": "new-msg-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				m := NewLiveChatMessage(opts...)
				var buf bytes.Buffer
				if err := m.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatMessage.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveChatMessage_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete live chat message",
			opts: []Option{
				WithIds([]string{"msg-id"}),
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
				m := NewLiveChatMessage(opts...)
				var buf bytes.Buffer
				if err := m.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatMessage.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveChatMessage_Transition(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "transition live chat message",
			opts: []Option{
				WithIds([]string{"msg-id"}),
				WithStatus("closed"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				q := r.URL.Query()
				if q.Get("id") != "msg-id" {
					t.Errorf("expected id=msg-id, got %s", q.Get("id"))
				}
				if q.Get("status") != "closed" {
					t.Errorf("expected status=closed, got %s", q.Get("status"))
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
							_, _ = w.Write([]byte(`{"id": "msg-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				m := NewLiveChatMessage(opts...)
				var buf bytes.Buffer
				if err := m.Transition(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveChatMessage.Transition() error = %v, wantErr %v", err,
						tt.wantErr,
					)
				}
			},
		)
	}
}
