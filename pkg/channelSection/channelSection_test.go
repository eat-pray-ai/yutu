package channelSection

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewChannelSection(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false

	tests := []struct {
		name string
		args args
		want ChannelSection
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"section1", "section2"}),
					WithChannelId("channel123"),
					WithHl("en"),
					WithMine(&mineTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithService(&youtube.Service{}),
				},
			},
			want: &channelSection{
				IDs:                    []string{"section1", "section2"},
				ChannelId:              "channel123",
				Hl:                     "en",
				Mine:                   &mineTrue,
				OnBehalfOfContentOwner: "owner123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &channelSection{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
				},
			},
			want: &channelSection{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
				},
			},
			want: &channelSection{
				Mine: &mineFalse,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithHl(""),
					WithOnBehalfOfContentOwner(""),
				},
			},
			want: &channelSection{
				ChannelId:              "",
				Hl:                     "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIDs([]string{"section1"}),
					WithChannelId("partialChannel"),
					WithHl("fr"),
				},
			},
			want: &channelSection{
				IDs:       []string{"section1"},
				ChannelId: "partialChannel",
				Hl:        "fr",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannelSection(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewChannelSection() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
