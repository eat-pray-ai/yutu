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
					WithMine(true, true),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &channelSection{
				ChannelId:              "channelId",
				Hl:                     "hl",
				Mine:                   &[]bool{true}[0],
				OnBehalfOfContentOwner: "contentOwner",
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
