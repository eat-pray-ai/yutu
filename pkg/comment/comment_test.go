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

	tests := []struct {
		name string
		args args
		want Comment[youtube.Comment]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"comment1", "comment2"}),
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
					WithService(&youtube.Service{}),
				},
			},
			want: &comment{
				IDs:              []string{"comment1", "comment2"},
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
			want: &comment{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithCanRate(nil),
					WithBanAuthor(nil),
				},
			},
			want: &comment{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithCanRate(&canRateFalse),
					WithBanAuthor(&banAuthorFalse),
				},
			},
			want: &comment{
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
			want: &comment{
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
			want: &comment{
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
			want: &comment{
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
					WithIDs([]string{"comment1"}),
					WithTextOriginal("Partial comment"),
					WithVideoId("video456"),
					WithMaxResults(25),
				},
			},
			want: &comment{
				IDs:          []string{"comment1"},
				TextOriginal: "Partial comment",
				VideoId:      "video456",
				MaxResults:   25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewComment(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewComment() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
