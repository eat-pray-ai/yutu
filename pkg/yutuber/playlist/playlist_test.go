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
					WithID("id"),
					WithTitle("title"),
					WithDescription("description"),
					WithHl("hl"),
					WithMaxResults(5),
					WithMine(true, true),
					WithTags([]string{"tag1", "tag2"}),
					WithLanguage("language"),
					WithChannelId("channelId"),
					WithPrivacy("public"),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOnBehalfOfContentOwnerChannel("contentOwnerChannel"),
				},
			},
			want: &playlist{
				ID:                            "id",
				Title:                         "title",
				Description:                   "description",
				Hl:                            "hl",
				MaxResults:                    5,
				Mine:                          &[]bool{true}[0],
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
