package i18nRegion

import (
	"reflect"
	"testing"
)

func TestNewI18nRegion(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want I18nRegion
	}{
		{
			name: "TestNewI18nRegion",
			args: args{
				opts: []Option{
					WithHl("hl"),
				},
			},
			want: &i18nRegion{
				hl: "hl",
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
