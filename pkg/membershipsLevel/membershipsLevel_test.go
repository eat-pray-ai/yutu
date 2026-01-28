// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewMembershipsLevel(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IMembershipsLevel[youtube.MembershipsLevel]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("items"),
					WithService(svc),
				},
			},
			want: &MembershipsLevel{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "items",
				},
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &MembershipsLevel{Fields: &common.Fields{}},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewMembershipsLevel(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewMembershipsLevel() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}
