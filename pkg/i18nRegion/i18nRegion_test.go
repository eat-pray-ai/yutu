// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewI18nRegion(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want II18nRegion[youtube.I18nRegion]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[*].id"),
					WithService(svc),
				},
			},
			want: &I18nRegion{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[*].id",
				},
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &I18nRegion{Fields: &common.Fields{}},
		},
		{
			name: "with empty string value",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &I18nRegion{
				Fields: &common.Fields{},
				Hl:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewI18nRegion(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewI18nRegion() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
