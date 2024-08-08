package channelBanner

import (
	"reflect"
	"testing"
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
			name: "TestNewChannelBanner",
			args: args{
				opts: []Option{
					WithFile("file"),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOnBehalfOfContentOwnerChannel("contentOwnerChannel"),
				},
			},
			want: &channelBanner{
				File:                          "file",
				OnBehalfOfContentOwner:        "contentOwner",
				OnBehalfOfContentOwnerChannel: "contentOwnerChannel",
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
