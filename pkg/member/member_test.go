// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewMember(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IMember[youtube.Member]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithMemberChannelId("member123"),
					WithHasAccessToLevel("level1"),
					WithMaxResults(100),
					WithMode("all_current"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &Member{
				DefaultFields: &pkg.DefaultFields{
					Service:  svc,
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				MemberChannelId:  "member123",
				HasAccessToLevel: "level1",
				MaxResults:       100,
				Mode:             "all_current",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Member{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Member{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-15),
				},
			},
			want: &Member{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithMemberChannelId(""),
					WithHasAccessToLevel(""),
					WithMode(""),
				},
			},
			want: &Member{
				DefaultFields:    &pkg.DefaultFields{},
				MemberChannelId:  "",
				HasAccessToLevel: "",
				Mode:             "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithMemberChannelId("channel456"),
					WithMaxResults(50),
					WithParts([]string{"id"}),
				},
			},
			want: &Member{
				DefaultFields:   &pkg.DefaultFields{Parts: []string{"id"}},
				MemberChannelId: "channel456",
				MaxResults:      50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewMember(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewMember() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
