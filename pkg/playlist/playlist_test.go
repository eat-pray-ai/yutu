package playlist

import (
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"reflect"
	"testing"
)

func TestNewPlaylist(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Playlist
	}{
		{
			name: "TestNewPlaylist",
			args: args{
				opts: []Option{
					WithIDs([]string{"id1", "id2"}),
					WithTitle("title"),
					WithDescription("description"),
					WithHl("hl"),
					WithMaxResults(5),
					WithMine(utils.BoolPtr("true")),
					WithTags([]string{"tag1", "tag2"}),
					WithLanguage("language"),
					WithChannelId("channelId"),
					WithPrivacy("public"),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOnBehalfOfContentOwnerChannel("contentOwnerChannel"),
				},
			},
			want: &playlist{
				IDs:                           []string{"id1", "id2"},
				Title:                         "title",
				Description:                   "description",
				Hl:                            "hl",
				MaxResults:                    5,
				Mine:                          utils.BoolPtr("true"),
				Tags:                          []string{"tag1", "tag2"},
				Language:                      "language",
				ChannelId:                     "channelId",
				Privacy:                       "public",
				OnBehalfOfContentOwner:        "contentOwner",
				OnBehalfOfContentOwnerChannel: "contentOwnerChannel",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylist(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewPlaylist() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
