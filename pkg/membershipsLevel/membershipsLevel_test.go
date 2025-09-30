// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewMembershipsLevel(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want MembershipsLevel[youtube.MembershipsLevel]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithService(&youtube.Service{}),
				},
			},
			want: &membershipsLevel{},
		},
		{
			name: "with no options",
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
