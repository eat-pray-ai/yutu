// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package commentThread

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewCommentThread(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want CommentThread[youtube.CommentThread]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"thread1", "thread2"}),
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
					WithService(&youtube.Service{}),
				},
			},
			want: &commentThread{
				IDs:                          []string{"thread1", "thread2"},
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
			want: &commentThread{},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &commentThread{
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
			want: &commentThread{
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
				},
			},
			want: &commentThread{
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
					WithIDs([]string{"thread1"}),
					WithVideoId("video456"),
					WithTextOriginal("Partial comment thread"),
					WithMaxResults(50),
				},
			},
			want: &commentThread{
				IDs:          []string{"thread1"},
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
					t.Errorf("NewCommentThread() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
