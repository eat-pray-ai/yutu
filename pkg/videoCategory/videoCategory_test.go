// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewVideoCategory(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IVideoCategory[youtube.VideoCategory]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"cat1", "cat2"}),
					WithHl("en"),
					WithRegionCode("US"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &VideoCategory{
				Ids:        []string{"cat1", "cat2"},
				Hl:         "en",
				RegionCode: "US",
				Parts:      []string{"snippet"},
				Output:     "json",
				Jsonpath:   "items.id",
				service:    svc,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &VideoCategory{},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
					WithRegionCode(""),
				},
			},
			want: &VideoCategory{
				Hl:         "",
				RegionCode: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHl("ja"),
					WithRegionCode("JP"),
				},
			},
			want: &VideoCategory{
				Hl:         "ja",
				RegionCode: "JP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoCategory(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewVideoCategory() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
