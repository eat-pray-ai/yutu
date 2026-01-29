// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

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

func TestNewComment(t *testing.T) {
	type args struct {
		opts []Option
	}

	canRateTrue := true
	canRateFalse := false
	banAuthorTrue := true
	banAuthorFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IComment[youtube.Comment]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"comment1", "comment2"}),
					WithAuthorChannelId("author123"),
					WithCanRate(&canRateTrue),
					WithChannelId("channel123"),
					WithMaxResults(50),
					WithParentId("parent123"),
					WithTextFormat("html"),
					WithTextOriginal("This is a comment"),
					WithModerationStatus("published"),
					WithBanAuthor(&banAuthorTrue),
					WithVideoId("video123"),
					WithViewerRating("like"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &Comment{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				Ids:              []string{"comment1", "comment2"},
				AuthorChannelId:  "author123",
				CanRate:          &canRateTrue,
				ChannelId:        "channel123",
				MaxResults:       50,
				ParentId:         "parent123",
				TextFormat:       "html",
				TextOriginal:     "This is a comment",
				ModerationStatus: "published",
				BanAuthor:        &banAuthorTrue,
				VideoId:          "video123",
				ViewerRating:     "like",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Comment{Fields: &common.Fields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithCanRate(nil),
					WithBanAuthor(nil),
				},
			},
			want: &Comment{Fields: &common.Fields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithCanRate(&canRateFalse),
					WithBanAuthor(&banAuthorFalse),
				},
			},
			want: &Comment{
				Fields:    &common.Fields{},
				CanRate:   &canRateFalse,
				BanAuthor: &banAuthorFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Comment{
				Fields:     &common.Fields{},
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &Comment{
				Fields:     &common.Fields{},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithAuthorChannelId(""),
					WithChannelId(""),
					WithParentId(""),
					WithTextFormat(""),
					WithTextOriginal(""),
					WithModerationStatus(""),
					WithVideoId(""),
					WithViewerRating(""),
				},
			},
			want: &Comment{
				Fields:           &common.Fields{},
				AuthorChannelId:  "",
				ChannelId:        "",
				ParentId:         "",
				TextFormat:       "",
				TextOriginal:     "",
				ModerationStatus: "",
				VideoId:          "",
				ViewerRating:     "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"comment1"}),
					WithTextOriginal("Partial comment"),
					WithVideoId("video456"),
					WithMaxResults(25),
					WithService(svc),
				},
			},
			want: &Comment{
				Fields:       &common.Fields{Service: svc},
				Ids:          []string{"comment1"},
				TextOriginal: "Partial comment",
				VideoId:      "video456",
				MaxResults:   25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewComment(tt.args.opts...)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewComment() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestComment_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get comments by id",
			opts: []Option{
				WithIds([]string{"comment-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "comment-id" {
					t.Errorf("expected id=comment-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get comments by parentId",
			opts: []Option{
				WithParentId("parent-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("parentId") != "parent-id" {
					t.Errorf("expected parentId=parent-id, got %s", r.URL.Query().Get("parentId"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get comments with textFormat",
			opts: []Option{
				WithTextFormat("plainText"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("textFormat") != "plainText" {
					t.Errorf("expected textFormat=plainText, got %s", r.URL.Query().Get("textFormat"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{
					"items": [
						{"id": "comment-1", "snippet": {"textOriginal": "Comment 1"}}
					]
				}`))
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
			c := NewComment(opts...)
			got, err := c.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Comment.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Comment.Get() got length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestComment_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := 0; i < 20; i++ {
				items[i] = fmt.Sprintf(`{"id": "comment-%d"}`, i)
			}
			w.Write([]byte(fmt.Sprintf(`{
				"items": [%s],
				"nextPageToken": "page-2"
			}`, strings.Join(items, ","))))
		} else if pageToken == "page-2" {
			w.Write([]byte(`{
				"items": [{"id": "comment-20"}, {"id": "comment-21"}],
				"nextPageToken": ""
			}`))
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

	c := NewComment(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := c.Get()
	if err != nil {
		t.Errorf("Comment.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Comment.Get() got length = %v, want 22", len(got))
	}
}

func TestComment_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"id": "comment-1",
					"snippet": {
						"authorDisplayName": "User",
						"videoId": "video-1",
						"textDisplay": "Comment"
					}
				}
			]
		}`))
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

	tests := []struct {
		name    string
		opts    []Option
		output  string
		wantErr bool
	}{
		{
			name: "list comments json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"comment-1"}),
				WithMaxResults(1),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list comments yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"comment-1"}),
				WithMaxResults(1),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list comments table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"comment-1"}),
				WithMaxResults(1),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewComment(tt.opts...)
			var buf bytes.Buffer
			if err := c.List(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Comment.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if buf.Len() == 0 {
				t.Errorf("Comment.List() output is empty")
			}
		})
	}
}

func TestComment_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert comment",
			opts: []Option{
				WithParentId("parent-id"),
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
			name: "insert comment full",
			opts: []Option{
				WithParentId("parent-id"),
				WithTextOriginal("New comment"),
				WithAuthorChannelId("author-id"),
				WithChannelId("channel-id"),
				WithVideoId("video-id"),
				func(c *Comment) {
					b := true
					c.CanRate = &b
				},
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
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"id": "new-comment-id"}`))
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
			c := NewComment(opts...)
			var buf bytes.Buffer
			if err := c.Insert(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Comment.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComment_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update comment",
			opts: []Option{
				WithIds([]string{"comment-id"}),
				WithTextOriginal("Updated text"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "snippet" {
						t.Errorf("expected part=snippet, got %s", r.URL.Query().Get("part"))
					}
				} else if r.Method == "GET" {
				} else {
					t.Errorf("unexpected method %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "update comment with canRate and viewerRating",
			opts: []Option{
				WithIds([]string{"comment-id"}),
				func(c *Comment) {
					b := true
					c.CanRate = &b
				},
				WithViewerRating("like"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
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
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "GET" {
					w.Write([]byte(`{
						"items": [
							{"id": "comment-id", "snippet": {"textOriginal": "Old Text"}}
						]
					}`))
				} else {
					w.Write([]byte(`{"id": "comment-id", "snippet": {"textOriginal": "Updated Text"}}`))
				}
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
			c := NewComment(opts...)
			var buf bytes.Buffer
			if err := c.Update(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Comment.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComment_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete comment",
			opts: []Option{
				WithIds([]string{"comment-id"}),
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
			c := NewComment(opts...)
			var buf bytes.Buffer
			if err := c.Delete(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Comment.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComment_MarkAsSpam(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "mark comment as spam",
			opts: []Option{
				WithIds([]string{"comment-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("id") != "comment-id" {
					t.Errorf("expected id=comment-id, got %s", r.URL.Query().Get("id"))
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
			c := NewComment(opts...)
			var buf bytes.Buffer
			if err := c.MarkAsSpam(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Comment.MarkAsSpam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestComment_SetModerationStatus(t *testing.T) {
	banAuthorTrue := true
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "set moderation status",
			opts: []Option{
				WithIds([]string{"comment-id"}),
				WithModerationStatus("published"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("id") != "comment-id" {
					t.Errorf("expected id=comment-id, got %s", r.URL.Query().Get("id"))
				}
				if r.URL.Query().Get("moderationStatus") != "published" {
					t.Errorf("expected moderationStatus=published, got %s", r.URL.Query().Get("moderationStatus"))
				}
			},
			wantErr: false,
		},
		{
			name: "set moderation status with banAuthor",
			opts: []Option{
				WithIds([]string{"comment-id"}),
				WithModerationStatus("published"),
				WithBanAuthor(&banAuthorTrue),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("banAuthor") != "true" {
					t.Errorf("expected banAuthor=true, got %s", r.URL.Query().Get("banAuthor"))
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
			c := NewComment(opts...)
			var buf bytes.Buffer
			if err := c.SetModerationStatus(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Comment.SetModerationStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
