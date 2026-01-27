// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewThumbnail(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want IThumbnail
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithVideoId("video123"),
					WithFile("/path/to/thumbnail.jpg"),
					WithOutput("json"),
					WithJsonpath("id"),
					WithService(svc),
				},
			},
			want: &Thumbnail{
				DefaultFields: &pkg.DefaultFields{
					Service:  svc,
					Output:   "json",
					Jsonpath: "id",
				},
				VideoId: "video123",
				File:    "/path/to/thumbnail.jpg",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Thumbnail{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithVideoId(""),
					WithFile(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &Thumbnail{
				DefaultFields: &pkg.DefaultFields{
					Output:   "",
					Jsonpath: "",
				},
				VideoId: "",
				File:    "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithVideoId("myVideo123"),
					WithFile("/images/thumb.png"),
				},
			},
			want: &Thumbnail{
				DefaultFields: &pkg.DefaultFields{},
				VideoId:       "myVideo123",
				File:          "/images/thumb.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewThumbnail(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewThumbnail() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
