package playlistItem

import (
	"reflect"
	"testing"
)

func TestNewPlaylistItem(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want PlaylistItem
	}{
		{
			name: "TestNewPlaylistItem",
			args: args{
				opts: []Option{
					WithId("id"),
					WithTitle("title"),
					WithDescription("description"),
					WithKind("video"),
					WithKVideoId("kVideoId"),
					WithPlaylistId("playlistId"),
					WithChannelId("channelId"),
					WithPrivacy("public"),
				},
			},
			want: &playlistItem{
				id:          "id",
				title:       "title",
				description: "description",
				kind:        "video",
				kVideoId:    "kVideoId",
				playlistId:  "playlistId",
				channelId:   "channelId",
				privacy:     "public",
			},
		},
		{
			name: "TestNewPlaylistItem",
			args: args{
				opts: []Option{
					WithId("id"),
					WithTitle("title"),
					WithDescription("description"),
					WithKind("channel"),
					WithKChannelId("kChannelId"),
					WithPlaylistId("playlistId"),
					WithChannelId("channelId"),
					WithPrivacy("private"),
				},
			},
			want: &playlistItem{
				id:          "id",
				title:       "title",
				description: "description",
				kind:        "channel",
				kChannelId:  "kChannelId",
				playlistId:  "playlistId",
				channelId:   "channelId",
				privacy:     "private",
			},
		},
		{
			name: "TestNewPlaylistItem",
			args: args{
				opts: []Option{
					WithVideoId("videoId"),
					WithMaxResults(5),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &playlistItem{
				videoId:                "videoId",
				maxResults:             5,
				onBehalfOfContentOwner: "contentOwner",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylistItem(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewPlaylistItem() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
