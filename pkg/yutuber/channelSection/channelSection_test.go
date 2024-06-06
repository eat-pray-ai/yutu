package channelSection

import (
	"reflect"
	"testing"
)

func TestNewChannelSection(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want ChannelSection
	}{
		{
			name: "TestNewChannelSection",
			args: args{
				opts: []Option{
					WithChannelId("channelId"),
					WithHl("hl"),
					WithMine("true"),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &channelSection{
				channelId:              "channelId",
				hl:                     "hl",
				mine:                   "true",
				onBehalfOfContentOwner: "contentOwner",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannelSection(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewChannelSection() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
