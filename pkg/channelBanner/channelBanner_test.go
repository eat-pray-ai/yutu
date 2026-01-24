// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelBanner

import (
	"reflect"
	"testing"

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
				ChannelId:                     "channel123",
				File:                          "/path/to/banner.jpg",
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
				Output:                        "json",
				Jsonpath:                      "items.id",
				service:                       svc,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &ChannelBanner{},
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
				ChannelId:                     "",
				File:                          "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
				Output:                        "",
				Jsonpath:                      "",
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
				ChannelId: "partialChannel",
				File:      "/partial/banner.png",
				Output:    "yaml",
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
