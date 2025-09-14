package i18nRegion

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewI18nRegion(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want I18nRegion[youtube.I18nRegion]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithService(&youtube.Service{}),
				},
			},
			want: &i18nRegion{
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &i18nRegion{},
		},
		{
			name: "with empty string value",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &i18nRegion{
				Hl: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewI18nRegion(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewI18nRegion() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
