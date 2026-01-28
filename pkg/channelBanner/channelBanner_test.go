// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewChannelBanner(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want IChannelBanner
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("channel123"),
					WithFile("/path/to/banner.jpg"),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &ChannelBanner{
				Fields: &common.Fields{
					Output:   "json",
					Jsonpath: "items.id",
					Service:  svc,
				},
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
			want: &ChannelBanner{Fields: &common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithFile(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &ChannelBanner{
				Fields: &common.Fields{
					Output:   "",
					Jsonpath: "",
				},
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
					WithOutput("yaml"),
				},
			},
			want: &ChannelBanner{
				Fields:    &common.Fields{Output: "yaml"},
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
					t.Errorf("%s\nNewChannelBanner() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
