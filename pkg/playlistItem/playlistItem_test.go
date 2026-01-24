// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewPlaylistItem(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want PlaylistItem[youtube.PlaylistItem]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"item1", "item2"}),
					WithTitle("Test Item"),
					WithDescription("Test item description"),
					WithKind("video"),
					WithKVideoId("video123"),
					WithKChannelId("channel123"),
					WithKPlaylistId("playlist123"),
					WithVideoId("video456"),
					WithPlaylistId("playlist456"),
					WithChannelId("channel456"),
					WithPrivacy("public"),
					WithMaxResults(50),
					WithOnBehalfOfContentOwner("owner123"),
					WithService(&youtube.Service{}),
				},
			},
			want: &playlistItem{
				Ids:                    []string{"item1", "item2"},
				Title:                  "Test Item",
				Description:            "Test item description",
				Kind:                   "video",
				KVideoId:               "video123",
				KChannelId:             "channel123",
				KPlaylistId:            "playlist123",
				VideoId:                "video456",
				PlaylistId:             "playlist456",
				ChannelId:              "channel456",
				Privacy:                "public",
				MaxResults:             50,
				OnBehalfOfContentOwner: "owner123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &playlistItem{},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &playlistItem{
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-15),
				},
			},
			want: &playlistItem{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithTitle(""),
					WithDescription(""),
					WithKind(""),
					WithKVideoId(""),
					WithKChannelId(""),
					WithKPlaylistId(""),
					WithVideoId(""),
					WithPlaylistId(""),
					WithChannelId(""),
					WithPrivacy(""),
					WithOnBehalfOfContentOwner(""),
				},
			},
			want: &playlistItem{
				Title:                  "",
				Description:            "",
				Kind:                   "",
				KVideoId:               "",
				KChannelId:             "",
				KPlaylistId:            "",
				VideoId:                "",
				PlaylistId:             "",
				ChannelId:              "",
				Privacy:                "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Video"),
					WithKind("video"),
					WithKVideoId("myVideo123"),
					WithPlaylistId("myPlaylist"),
					WithMaxResults(25),
				},
			},
			want: &playlistItem{
				Title:      "My Video",
				Kind:       "video",
				KVideoId:   "myVideo123",
				PlaylistId: "myPlaylist",
				MaxResults: 25,
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
