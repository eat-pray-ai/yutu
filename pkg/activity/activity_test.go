// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewActivity(t *testing.T) {
	type args struct {
		opts []Option
	}

	homeTrue := true
	homeFalse := false
	mineTrue := true
	mineFalse := false

	tests := []struct {
		name string
		args args
		want IActivity[youtube.Activity]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("test-channel-123"),
					WithHome(&homeTrue),
					WithMaxResults(50),
					WithMine(&mineTrue),
					WithPublishedAfter("2024-01-01T00:00:00Z"),
					WithPublishedBefore("2024-12-31T23:59:59Z"),
					WithRegionCode("US"),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{
				ChannelId:       "test-channel-123",
				Home:            &homeTrue,
				MaxResults:      50,
				Mine:            &mineTrue,
				PublishedAfter:  "2024-01-01T00:00:00Z",
				PublishedBefore: "2024-12-31T23:59:59Z",
				RegionCode:      "US",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithHome(nil),
					WithMine(nil),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithHome(&homeFalse),
					WithMine(&mineFalse),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{
				Home: &homeFalse,
				Mine: &mineFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithPublishedAfter(""),
					WithPublishedBefore(""),
					WithRegionCode(""),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{
				ChannelId:       "",
				PublishedAfter:  "",
				PublishedBefore: "",
				RegionCode:      "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("partial-channel"),
					WithMaxResults(25),
					WithRegionCode("UK"),
					WithService(&youtube.Service{}),
				},
			},
			want: &Activity{
				ChannelId:  "partial-channel",
				MaxResults: 25,
				RegionCode: "UK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewActivity(tt.args.opts...); !reflect.DeepEqual(got.(*Activity).ChannelId, tt.want.(*Activity).ChannelId) ||
					!reflect.DeepEqual(got.(*Activity).Home, tt.want.(*Activity).Home) ||
					!reflect.DeepEqual(got.(*Activity).MaxResults, tt.want.(*Activity).MaxResults) ||
					!reflect.DeepEqual(got.(*Activity).Mine, tt.want.(*Activity).Mine) ||
					!reflect.DeepEqual(got.(*Activity).PublishedAfter, tt.want.(*Activity).PublishedAfter) ||
					!reflect.DeepEqual(got.(*Activity).PublishedBefore, tt.want.(*Activity).PublishedBefore) ||
					!reflect.DeepEqual(got.(*Activity).RegionCode, tt.want.(*Activity).RegionCode) {
					t.Errorf("NewActivity() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
