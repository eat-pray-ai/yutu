package thumbnail

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewThumbnail(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want Thumbnail
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithVideoId("video123"),
					WithFile("/path/to/thumbnail.jpg"),
					WithService(&youtube.Service{}),
				},
			},
			want: &thumbnail{
				VideoId: "video123",
				File:    "/path/to/thumbnail.jpg",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &thumbnail{},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithVideoId(""),
					WithFile(""),
				},
			},
			want: &thumbnail{
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
			want: &thumbnail{
				VideoId: "myVideo123",
				File:    "/images/thumb.png",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewThumbnail(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewThumbnail() = %v, want %v", got, tt.want)
			}
		})
	}
}
