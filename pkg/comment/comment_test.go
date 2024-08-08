package comment

import (
	"reflect"
	"testing"
)

func TestNewComment(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Comment
	}{
		{
			name: "TestNewComment",
			args: args{
				opts: []Option{
					WithIDs([]string{"id1", "id2"}),
					WithAuthorChannelId("authorChannelId"),
					WithCanRate(true, true),
					WithChannelId("channelId"),
					WithMaxResults(5),
					WithParentId("parentId"),
					WithTextFormat("html"),
					WithTextOriginal("textOriginal"),
					WithModerationStatus("heldForReview"),
					WithBanAuthor(true, true),
					WithVideoId("videoId"),
					WithViewerRating("like"),
				},
			},
			want: &comment{
				IDs:              []string{"id1", "id2"},
				AuthorChannelId:  "authorChannelId",
				CanRate:          &[]bool{true}[0],
				ChannelId:        "channelId",
				MaxResults:       5,
				ParentId:         "parentId",
				TextFormat:       "html",
				TextOriginal:     "textOriginal",
				ModerationStatus: "heldForReview",
				BanAuthor:        &[]bool{true}[0],
				VideoId:          "videoId",
				ViewerRating:     "like",
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
