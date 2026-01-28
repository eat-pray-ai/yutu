// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package watermark

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewWatermark(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IWatermark
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("channel123"),
					WithFile("/path/to/watermark.png"),
					WithInVideoPosition("topRight"),
					WithDurationMs(5000),
					WithOffsetMs(1000),
					WithOffsetType("offsetFromStart"),
					WithOnBehalfOfContentOwner("owner123"),
					WithService(svc),
				},
			},
			want: &Watermark{
				Fields:                 &common.Fields{Service: svc},
				ChannelId:              "channel123",
				File:                   "/path/to/watermark.png",
				InVideoPosition:        "topRight",
				DurationMs:             5000,
				OffsetMs:               1000,
				OffsetType:             "offsetFromStart",
				OnBehalfOfContentOwner: "owner123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Watermark{Fields: &common.Fields{}},
		},
		{
			name: "with zero values",
			args: args{
				opts: []Option{
					WithDurationMs(0),
					WithOffsetMs(0),
				},
			},
			want: &Watermark{
				Fields:     &common.Fields{},
				DurationMs: 0,
				OffsetMs:   0,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithFile(""),
					WithInVideoPosition(""),
					WithOffsetType(""),
					WithOnBehalfOfContentOwner(""),
				},
			},
			want: &Watermark{
				Fields:                 &common.Fields{},
				ChannelId:              "",
				File:                   "",
				InVideoPosition:        "",
				OffsetType:             "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("myChannel"),
					WithFile("/watermarks/logo.png"),
					WithInVideoPosition("bottomLeft"),
					WithDurationMs(10000),
				},
			},
			want: &Watermark{
				Fields:          &common.Fields{},
				ChannelId:       "myChannel",
				File:            "/watermarks/logo.png",
				InVideoPosition: "bottomLeft",
				DurationMs:      10000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewWatermark(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewWatermark() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
