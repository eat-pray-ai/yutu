package video

import (
	"reflect"
	"testing"
)

func TestNewVideo(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Video
	}{
		{
			name: "TestNewVideo",
			args: args{
				opts: []Option{
					WithId("id"),
					WithAutoLevels(true, true),
					WithFile("file"),
					WithTitle("title"),
					WithDescription("description"),
					WithHl("hl"),
					WithTags([]string{"tag1", "tag2"}),
					WithLanguage("language"),
					WithLocale("locale"),
					WithLicense("license"),
				},
			},
			want: &video{
				id:          "id",
				autoLevels:  &[]bool{true}[0],
				file:        "file",
				title:       "title",
				description: "description",
				hl:          "hl",
				tags:        []string{"tag1", "tag2"},
				language:    "language",
				locale:      "locale",
				license:     "license",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideo(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewVideo() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
