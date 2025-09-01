package channelBanner

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewChannelBanner(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want ChannelBanner
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("channel123"),
					WithFile("/path/to/banner.jpg"),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithService(&youtube.Service{}),
				},
			},
			want: &channelBanner{
				ChannelId:                     "channel123",
				File:                          "/path/to/banner.jpg",
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &channelBanner{},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithFile(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &channelBanner{
				ChannelId:                     "",
				File:                          "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("partialChannel"),
					WithFile("/partial/banner.png"),
				},
			},
			want: &channelBanner{
				ChannelId: "partialChannel",
				File:      "/partial/banner.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannelBanner(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewChannelBanner() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
