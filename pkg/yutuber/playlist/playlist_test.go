package playlist

import (
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
					WithId("id"),
					WithTitle("title"),
					WithDescription("description"),
					WithHl("hl"),
					WithMaxResults(5),
					WithMine("true"),
					WithTags([]string{"tag1", "tag2"}),
					WithLanguage("language"),
					WithChannelId("channelId"),
					WithPrivacy("public"),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOnBehalfOfContentOwnerChannel("contentOwnerChannel"),
				},
			},
			want: &playlist{
				id:                            "id",
				title:                         "title",
				description:                   "description",
				hl:                            "hl",
				maxResults:                    5,
				mine:                          "true",
				tags:                          []string{"tag1", "tag2"},
				language:                      "language",
				channelId:                     "channelId",
				privacy:                       "public",
				onBehalfOfContentOwner:        "contentOwner",
				onBehalfOfContentOwnerChannel: "contentOwnerChannel",
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
