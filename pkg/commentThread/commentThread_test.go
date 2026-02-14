// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestNewCommentThread(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want ICommentThread[youtube.CommentThread]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"thread1", "thread2"}),
					WithAllThreadsRelatedToChannelId("relatedChannel123"),
					WithAuthorChannelId("author123"),
					WithChannelId("channel123"),
					WithMaxResults(100),
					WithModerationStatus("published"),
					WithOrder("time"),
					WithSearchTerms("test search"),
					WithTextFormat("html"),
					WithTextOriginal("This is a comment thread"),
					WithVideoId("video123"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &CommentThread{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				Ids:                          []string{"thread1", "thread2"},
				AllThreadsRelatedToChannelId: "relatedChannel123",
				AuthorChannelId:              "author123",
				ChannelId:                    "channel123",
				MaxResults:                   100,
				ModerationStatus:             "published",
				Order:                        "time",
				SearchTerms:                  "test search",
				TextFormat:                   "html",
				TextOriginal:                 "This is a comment thread",
				VideoId:                      "video123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &CommentThread{Fields: &common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &CommentThread{
				Fields:     &common.Fields{},
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-20),
				},
			},
			want: &CommentThread{
				Fields:     &common.Fields{},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithAllThreadsRelatedToChannelId(""),
					WithAuthorChannelId(""),
					WithChannelId(""),
					WithModerationStatus(""),
					WithOrder(""),
					WithSearchTerms(""),
					WithTextFormat(""),
					WithTextOriginal(""),
					WithVideoId(""),
					WithParts(nil),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &CommentThread{
				Fields: &common.Fields{
					Parts:    nil,
					Output:   "",
					Jsonpath: "",
				},
				AllThreadsRelatedToChannelId: "",
				AuthorChannelId:              "",
				ChannelId:                    "",
				ModerationStatus:             "",
				Order:                        "",
				SearchTerms:                  "",
				TextFormat:                   "",
				TextOriginal:                 "",
				VideoId:                      "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"thread1"}),
					WithVideoId("video456"),
					WithTextOriginal("Partial comment thread"),
					WithMaxResults(50),
					WithOutput("yaml"),
				},
			},
			want: &CommentThread{
				Fields:       &common.Fields{Output: "yaml"},
				Ids:          []string{"thread1"},
				VideoId:      "video456",
				TextOriginal: "Partial comment thread",
				MaxResults:   50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewCommentThread(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewCommentThread() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestCommentThread_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get comment threads by id",
			opts: []Option{
				WithIds([]string{"thread-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "thread-id" {
					t.Errorf("expected id=thread-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get comment threads by allThreadsRelatedToChannelId",
			opts: []Option{
				WithAllThreadsRelatedToChannelId("channel-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("allThreadsRelatedToChannelId") != "channel-id" {
					t.Errorf(
						"expected allThreadsRelatedToChannelId=channel-id, got %s",
						r.URL.Query().Get("allThreadsRelatedToChannelId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get comment threads by channelId",
			opts: []Option{
				WithChannelId("channel-id"),
				WithMaxResults(1),
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
			name: "get comment threads by videoId",
			opts: []Option{
				WithVideoId("video-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("videoId") != "video-id" {
					t.Errorf(
						"expected videoId=video-id, got %s", r.URL.Query().Get("videoId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get comment threads with options",
			opts: []Option{
				WithChannelId("channel-id"),
				WithModerationStatus("heldForReview"),
				WithOrder("relevance"),
				WithSearchTerms("search"),
				WithTextFormat("plainText"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("moderationStatus") != "heldForReview" {
					t.Errorf(
						"expected moderationStatus=heldForReview, got %s",
						q.Get("moderationStatus"),
					)
				}
				if q.Get("order") != "relevance" {
					t.Errorf("expected order=relevance, got %s", q.Get("order"))
				}
				if q.Get("searchTerms") != "search" {
					t.Errorf("expected searchTerms=search, got %s", q.Get("searchTerms"))
				}
				if q.Get("textFormat") != "plainText" {
					t.Errorf("expected textFormat=plainText, got %s", q.Get("textFormat"))
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
						{"id": "thread-1", "snippet": {"topLevelComment": {"snippet": {"textDisplay": "Comment 1"}}}}
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
				c := NewCommentThread(opts...)
				got, err := c.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("CommentThread.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"CommentThread.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestCommentThread_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := range 20 {
				items[i] = fmt.Sprintf(`{"id": "thread-%d"}`, i)
			}
			_, _ = w.Write(
				fmt.Appendf(nil,
					`{
				"items": [%s],
				"nextPageToken": "page-2"
			}`, strings.Join(items, ","),
				),
			)
		} else if pageToken == "page-2" {
			_, _ = w.Write(
				[]byte(`{
				"items": [{"id": "thread-20"}, {"id": "thread-21"}],
				"nextPageToken": ""
			}`),
			)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	svc, err := youtube.NewService(
		context.Background(),
		option.WithEndpoint(ts.URL),
		option.WithAPIKey("test-key"),
	)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	c := NewCommentThread(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := c.Get()
	if err != nil {
		t.Errorf("CommentThread.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("CommentThread.Get() got length = %v, want 22", len(got))
	}
}

func TestCommentThread_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "thread-1",
					"snippet": {
						"topLevelComment": {
							"snippet": {
								"authorDisplayName": "User",
								"videoId": "video-1",
								"textDisplay": "Comment"
							}
						}
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
			name: "list comment threads json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"thread-1"}),
				WithMaxResults(1),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list comment threads yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"thread-1"}),
				WithMaxResults(1),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list comment threads table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"thread-1"}),
				WithMaxResults(1),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := NewCommentThread(tt.opts...)
				var buf bytes.Buffer
				if err := c.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"CommentThread.List() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
				if buf.Len() == 0 {
					t.Errorf("CommentThread.List() output is empty")
				}
			},
		)
	}
}

func TestCommentThread_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert comment thread",
			opts: []Option{
				WithChannelId("channel-id"),
				WithVideoId("video-id"),
				WithTextOriginal("New comment"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert comment thread with author",
			opts: []Option{
				WithChannelId("channel-id"),
				WithVideoId("video-id"),
				WithTextOriginal("New comment"),
				WithAuthorChannelId("author-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
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
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"id": "new-thread-id"}`))
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
				c := NewCommentThread(opts...)
				var buf bytes.Buffer
				if err := c.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"CommentThread.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
