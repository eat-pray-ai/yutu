// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewSuperChatEvent(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ISuperChatEvent[youtube.SuperChatEvent]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithMaxResults(50),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[*].id"),
					WithService(svc),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[*].id",
				},
				Hl:         "en",
				MaxResults: 50,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &SuperChatEvent{Fields: &common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &SuperChatEvent{
				Fields:     &common.Fields{},
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &SuperChatEvent{
				Fields:     &common.Fields{},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
					WithOutput(""),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{Output: ""},
				Hl:     "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHl("ja"),
					WithMaxResults(25),
					WithOutput("yaml"),
				},
			},
			want: &SuperChatEvent{
				Fields:     &common.Fields{Output: "yaml"},
				Hl:         "ja",
				MaxResults: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSuperChatEvent(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewSuperChatEvent() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}
