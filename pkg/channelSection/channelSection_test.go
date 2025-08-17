package channelSection

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/utils"
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
					WithMine(utils.BoolPtr("true")),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &channelSection{
				ChannelId:              "channelId",
				Hl:                     "hl",
				Mine:                   utils.BoolPtr("true"),
				OnBehalfOfContentOwner: "contentOwner",
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
