package membershipsLevel

import (
	"reflect"
	"testing"
)

func TestNewMembershipsLevel(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want MembershipsLevel
	}{
		{
			name: "TestNewMembershipsLevel",
			args: args{
				opts: []Option{},
			},
			want: &membershipsLevel{},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewMembershipsLevel(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewMembershipsLevel() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
