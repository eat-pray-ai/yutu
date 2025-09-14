package playlist

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewPlaylist(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false

	tests := []struct {
		name string
		args args
		want Playlist
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"playlist1", "playlist2"}),
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
					WithService(&youtube.Service{}),
				},
			},
			want: &playlist{
				IDs:                           []string{"playlist1", "playlist2"},
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
			want: &playlist{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
				},
			},
			want: &playlist{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
				},
			},
			want: &playlist{
				Mine: &mineFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &playlist{
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &playlist{
				MaxResults: 1,
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
			want: &playlist{
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
			want: &playlist{
				Title:       "My Playlist",
				Description: "A great playlist",
				Privacy:     "private",
				MaxResults:  25,
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
