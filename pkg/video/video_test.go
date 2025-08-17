package video

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/utils"
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
					WithIDs([]string{"id1", "id2"}),
					WithAutoLevels(utils.BoolPtr("true")),
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
				IDs:         []string{"id1", "id2"},
				AutoLevels:  utils.BoolPtr("true"),
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
