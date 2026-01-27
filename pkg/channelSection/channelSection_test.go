// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channelSection

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewChannelSection(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IChannelSection[youtube.ChannelSection]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"section1", "section2"}),
					WithChannelId("channel123"),
					WithHl("en"),
					WithMine(&mineTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithService(svc),
				},
			},
			want: &ChannelSection{
				DefaultFields:          &pkg.DefaultFields{Service: svc},
				Ids:                    []string{"section1", "section2"},
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
			want: &ChannelSection{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
				},
			},
			want: &ChannelSection{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
				},
			},
			want: &ChannelSection{
				DefaultFields: &pkg.DefaultFields{},
				Mine:          &mineFalse,
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
			want: &ChannelSection{
				DefaultFields:          &pkg.DefaultFields{},
				ChannelId:              "",
				Hl:                     "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"section1"}),
					WithChannelId("partialChannel"),
					WithHl("fr"),
				},
			},
			want: &ChannelSection{
				DefaultFields: &pkg.DefaultFields{},
				Ids:           []string{"section1"},
				ChannelId:     "partialChannel",
				Hl:            "fr",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannelSection(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewChannelSection() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}
