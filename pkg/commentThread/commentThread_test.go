// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
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
					WithService(svc),
				},
			},
			want: &CommentThread{
				Fields: &common.Fields{
					Service:    svc,
					Parts:      []string{"id", "snippet"},
					Output:     "json",
					Ids:        []string{"thread1", "thread2"},
					MaxResults: 100,
					ChannelId:  "channel123",
				},
				AllThreadsRelatedToChannelId: "relatedChannel123",
				AuthorChannelId:              "author123",
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
				Fields: &common.Fields{MaxResults: math.MaxInt64},
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
				Fields: &common.Fields{MaxResults: 1},
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
				},
			},
			want: &CommentThread{
				Fields: &common.Fields{
					Parts:     nil,
					Output:    "",
					ChannelId: "",
				},
				AllThreadsRelatedToChannelId: "",
				AuthorChannelId:              "",
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
				Fields: &common.Fields{
					Output:     "yaml",
					Ids:        []string{"thread1"},
					MaxResults: 50,
				},
				VideoId:      "video456",
				TextOriginal: "Partial comment thread",
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
						{"id": "thread-1", "snippet": {"topLevelComment": {"snippet": {"textDisplay": "Comment 1"}}}}
					]
				}`),
							)
						},
					),
				)

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
	svc := common.NewTestService(t, common.PaginationHandler("thread"))
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
	common.RunListTest(
		t,
		`{
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
		}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			c := NewCommentThread(
				WithService(svc),
				WithOutput(output),
				WithIds([]string{"thread-1"}),
				WithMaxResults(1),
			)
			return c.List
		},
	)
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
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"id": "new-thread-id"}`))
						},
					),
				)

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
