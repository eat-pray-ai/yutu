// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package comment

import (
	"math"
	"reflect"
	"testing"

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
				Parts:            []string{"id", "snippet"},
				Output:           "json",
				Jsonpath:         "items.id",
				service:          svc,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Comment{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithCanRate(nil),
					WithBanAuthor(nil),
				},
			},
			want: &Comment{},
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
				Ids:          []string{"comment1"},
				TextOriginal: "Partial comment",
				VideoId:      "video456",
				MaxResults:   25,
				service:      svc,
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
