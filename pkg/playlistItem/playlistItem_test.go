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
					WithIDs([]string{"id1", "id2"}),
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
				IDs:         []string{"id1", "id2"},
				Title:       "title",
				Description: "description",
				Kind:        "video",
				KVideoId:    "kVideoId",
				PlaylistId:  "playlistId",
				ChannelId:   "channelId",
				Privacy:     "public",
			},
		},
		{
			name: "TestNewPlaylistItem",
			args: args{
				opts: []Option{
					WithIDs([]string{"id"}),
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
				IDs:         []string{"id"},
				Title:       "title",
				Description: "description",
				Kind:        "channel",
				KChannelId:  "kChannelId",
				PlaylistId:  "playlistId",
				ChannelId:   "channelId",
				Privacy:     "private",
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
				VideoId:                "videoId",
				MaxResults:             5,
				OnBehalfOfContentOwner: "contentOwner",
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
