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

	tests := []struct {
		name string
		args args
		want SuperChatEvent[youtube.SuperChatEvent]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithMaxResults(50),
					WithService(&youtube.Service{}),
				},
			},
			want: &superChatEvent{
				Hl:         "en",
				MaxResults: 50,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &superChatEvent{},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &superChatEvent{
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
			want: &superChatEvent{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &superChatEvent{
				Hl: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHl("ja"),
					WithMaxResults(25),
				},
			},
			want: &superChatEvent{
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
					t.Errorf("NewSuperChatEvent() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
