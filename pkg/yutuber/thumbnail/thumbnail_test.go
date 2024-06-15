package thumbnail

import (
	"reflect"
	"testing"
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
			name: "TestNewThumbnail",
			args: args{
				opts: []Option{
					WithFile("file"),
					WithVideoId("videoId"),
				},
			},
			want: &thumbnail{
				file:    "file",
				videoId: "videoId",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewThumbnail(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewThumbnail() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
