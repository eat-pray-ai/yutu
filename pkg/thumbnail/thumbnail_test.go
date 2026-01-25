// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thumbnail

import (
	"reflect"
	"testing"

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
				VideoId:  "video123",
				File:     "/path/to/thumbnail.jpg",
				Output:   "json",
				Jsonpath: "id",
				service:  svc,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Thumbnail{},
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
				VideoId:  "",
				File:     "",
				Output:   "",
				Jsonpath: "",
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
				VideoId: "myVideo123",
				File:    "/images/thumb.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewThumbnail(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s\nNewThumbnail() = %v\nwant %v", tt.name, got, tt.want)
			}
		})
	}
}
