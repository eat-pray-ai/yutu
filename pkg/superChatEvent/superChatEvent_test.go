// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

import (
	"math"
	"reflect"
	"testing"

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
				Hl:         "en",
				MaxResults: 50,
				Parts:      []string{"id", "snippet"},
				Output:     "json",
				Jsonpath:   "$.items[*].id",
				service:    svc,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &SuperChatEvent{},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &SuperChatEvent{MaxResults: math.MaxInt64},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &SuperChatEvent{MaxResults: 1},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
					WithOutput(""),
				},
			},
			want: &SuperChatEvent{Hl: "", Output: ""},
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
				Hl:         "ja",
				MaxResults: 25,
				Output:     "yaml",
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
