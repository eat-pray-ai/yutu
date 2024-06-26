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
					WithID("id"),
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
				ID:          "id",
				AutoLevels:  &[]bool{true}[0],
				File:        "file",
				Title:       "title",
				Description: "description",
				Hl:          "hl",
				Tags:        []string{"tag1", "tag2"},
				Language:    "language",
				Locale:      "locale",
				License:     "license",
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
