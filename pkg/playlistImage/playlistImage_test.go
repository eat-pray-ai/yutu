// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewPlaylistImage(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want IPlaylistImage[youtube.PlaylistImage]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"image1", "image2"}),
					WithHeight(1080),
					WithPlaylistId("playlist123"),
					WithType("hero"),
					WithWidth(1920),
					WithFile("/path/to/image.jpg"),
					WithParent("parent123"),
					WithMaxResults(50),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[*].id"),
					WithService(svc),
				},
			},
			want: &PlaylistImage{
				DefaultFields: &pkg.DefaultFields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[*].id",
				},
				Ids:                           []string{"image1", "image2"},
				Height:                        1080,
				PlaylistId:                    "playlist123",
				Type:                          "hero",
				Width:                         1920,
				File:                          "/path/to/image.jpg",
				Parent:                        "parent123",
				MaxResults:                    50,
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &PlaylistImage{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &PlaylistImage{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-20),
				},
			},
			want: &PlaylistImage{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithPlaylistId(""),
					WithType(""),
					WithFile(""),
					WithParent(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &PlaylistImage{
				DefaultFields:                 &pkg.DefaultFields{},
				PlaylistId:                    "",
				Type:                          "",
				File:                          "",
				Parent:                        "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithPlaylistId("myPlaylist"),
					WithType("hero"),
					WithFile("/images/hero.png"),
					WithMaxResults(25),
				},
			},
			want: &PlaylistImage{
				DefaultFields: &pkg.DefaultFields{},
				PlaylistId:    "myPlaylist",
				Type:          "hero",
				File:          "/images/hero.png",
				MaxResults:    25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylistImage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewPlaylistImage() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
