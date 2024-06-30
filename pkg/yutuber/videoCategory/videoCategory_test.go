package videoCategory

import (
	"reflect"
	"testing"
)

func TestNewVideoCategory(t *testing.T) {
	type args struct {
		opt []Option
	}
	tests := []struct {
		name string
		args args
		want VideoCategory
	}{
		{
			name: "TestNewVideoCategory",
			args: args{
				opt: []Option{
					WithID("id"),
					WithHl("hl"),
					WithRegionCode("regionCode"),
				},
			},
			want: &videoCategory{
				ID:         "id",
				Hl:         "hl",
				RegionCode: "regionCode",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoCategory(tt.args.opt...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewVideoCategory() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
