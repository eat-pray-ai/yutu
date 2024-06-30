package watermark

import (
	"reflect"
	"testing"
)

func TestNewWatermark(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Watermark
	}{
		{
			name: "TestNewWatermark",
			args: args{
				opts: []Option{
					WithChannelId("channelId"),
					WithFile("file"),
					WithInVideoPosition("topRight"),
					WithDurationMs(1024),
					WithOffsetMs(31415),
					WithOffsetType("offsetFromEnd"),
					WithOnBehalfOfContentOwner("contentOwner"),
				},
			},
			want: &watermark{
				ChannelId:              "channelId",
				File:                   "file",
				InVideoPosition:        "topRight",
				DurationMs:             1024,
				OffsetMs:               31415,
				OffsetType:             "offsetFromEnd",
				OnBehalfOfContentOwner: "contentOwner",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewWatermark(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewWatermark() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
