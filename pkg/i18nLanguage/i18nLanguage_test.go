package i18nLanguage

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewI18nLanguage(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want I18nLanguage[youtube.I18nLanguage]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithService(&youtube.Service{}),
				},
			},
			want: &i18nLanguage{
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &i18nLanguage{},
		},
		{
			name: "with empty string value",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &i18nLanguage{
				Hl: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewI18nLanguage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewI18nLanguage() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
