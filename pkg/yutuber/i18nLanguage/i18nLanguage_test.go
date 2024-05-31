package i18nLanguage

import (
	"reflect"
	"testing"
)

func TestNewI18nLanguage(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want I18nLanguage
	}{
		{
			name: "TestNewI18nLanguage",
			args: args{
				opts: []Option{
					WithHl("hl"),
				},
			},
			want: &i18nLanguage{
				hl: "hl",
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
