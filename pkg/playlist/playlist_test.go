// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewPlaylist(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IPlaylist[youtube.Playlist]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"playlist1", "playlist2"}),
					WithTitle("Test Playlist"),
					WithDescription("Test playlist description"),
					WithTags([]string{"tag1", "tag2", "tag3"}),
					WithLanguage("en"),
					WithChannelId("channel123"),
					WithPrivacy("public"),
					WithHl("en"),
					WithMaxResults(50),
					WithMine(&mineTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[0].snippet.title"),
					WithService(svc),
				},
			},
			want: &Playlist{
				DefaultFields: &pkg.DefaultFields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[0].snippet.title",
				},
				Ids:                           []string{"playlist1", "playlist2"},
				Title:                         "Test Playlist",
				Description:                   "Test playlist description",
				Tags:                          []string{"tag1", "tag2", "tag3"},
				Language:                      "en",
				ChannelId:                     "channel123",
				Privacy:                       "public",
				Hl:                            "en",
				MaxResults:                    50,
				Mine:                          &mineTrue,
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Playlist{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
				},
			},
			want: &Playlist{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
				},
			},
			want: &Playlist{
				DefaultFields: &pkg.DefaultFields{},
				Mine:          &mineFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Playlist{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &Playlist{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithTitle(""),
					WithDescription(""),
					WithLanguage(""),
					WithChannelId(""),
					WithPrivacy(""),
					WithHl(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &Playlist{
				DefaultFields:                 &pkg.DefaultFields{},
				Title:                         "",
				Description:                   "",
				Language:                      "",
				ChannelId:                     "",
				Privacy:                       "",
				Hl:                            "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Playlist"),
					WithDescription("A great playlist"),
					WithPrivacy("private"),
					WithMaxResults(25),
				},
			},
			want: &Playlist{
				DefaultFields: &pkg.DefaultFields{},
				Title:         "My Playlist",
				Description:   "A great playlist",
				Privacy:       "private",
				MaxResults:    25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylist(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewPlaylist() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
