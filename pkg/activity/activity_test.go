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
		want Activity[youtube.Activity]
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
			want: &activity{
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
			want: &activity{},
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
			want: &activity{},
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
			want: &activity{
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
			want: &activity{
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
			want: &activity{
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
			want: &activity{
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
			want: &activity{
				ChannelId:  "partial-channel",
				MaxResults: 25,
				RegionCode: "UK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewActivity(tt.args.opts...); !reflect.DeepEqual(got.(*activity).ChannelId, tt.want.(*activity).ChannelId) ||
					!reflect.DeepEqual(got.(*activity).Home, tt.want.(*activity).Home) ||
					!reflect.DeepEqual(got.(*activity).MaxResults, tt.want.(*activity).MaxResults) ||
					!reflect.DeepEqual(got.(*activity).Mine, tt.want.(*activity).Mine) ||
					!reflect.DeepEqual(got.(*activity).PublishedAfter, tt.want.(*activity).PublishedAfter) ||
					!reflect.DeepEqual(got.(*activity).PublishedBefore, tt.want.(*activity).PublishedBefore) ||
					!reflect.DeepEqual(got.(*activity).RegionCode, tt.want.(*activity).RegionCode) {
					t.Errorf("NewActivity() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
