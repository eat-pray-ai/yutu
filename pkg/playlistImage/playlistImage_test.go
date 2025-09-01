package playlistImage

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewPlaylistImage(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want PlaylistImage
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"image1", "image2"}),
					WithHeight(1080),
					WithPlaylistID("playlist123"),
					WithType("hero"),
					WithWidth(1920),
					WithFile("/path/to/image.jpg"),
					WithParent("parent123"),
					WithMaxResults(50),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithService(&youtube.Service{}),
				},
			},
			want: &playlistImage{
				IDs:                           []string{"image1", "image2"},
				Height:                        1080,
				PlaylistID:                    "playlist123",
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
			want: &playlistImage{},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &playlistImage{
				MaxResults: 1,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-20),
				},
			},
			want: &playlistImage{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithPlaylistID(""),
					WithType(""),
					WithFile(""),
					WithParent(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &playlistImage{
				PlaylistID:                    "",
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
					WithPlaylistID("myPlaylist"),
					WithType("hero"),
					WithFile("/images/hero.png"),
					WithMaxResults(25),
				},
			},
			want: &playlistImage{
				PlaylistID: "myPlaylist",
				Type:       "hero",
				File:       "/images/hero.png",
				MaxResults: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylistImage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewPlaylistImage() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
