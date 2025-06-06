package commentThread

import (
	"reflect"
	"testing"
)

func TestNewCommentThread(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want CommentThread
	}{
		{
			name: "TestNewCommentThread",
			args: args{
				opts: []Option{
					WithIDs([]string{"id"}),
					WithAllThreadsRelatedToChannelId("allThreadsRelatedToChannelId"),
					WithAuthorChannelId("authorChannelId"),
					WithChannelId("channelId"),
					WithMaxResults(5),
					WithModerationStatus("published"),
					WithOrder("relevance"),
					WithSearchTerms("searchTerms"),
					WithTextFormat("html"),
					WithTextOriginal("textOriginal"),
					WithVideoId("videoId"),
				},
			},
			want: &commentThread{
				IDs:                          []string{"id"},
				AllThreadsRelatedToChannelId: "allThreadsRelatedToChannelId",
				AuthorChannelId:              "authorChannelId",
				ChannelId:                    "channelId",
				MaxResults:                   5,
				ModerationStatus:             "published",
				Order:                        "relevance",
				SearchTerms:                  "searchTerms",
				TextFormat:                   "html",
				TextOriginal:                 "textOriginal",
				VideoId:                      "videoId",
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
