// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewI18nLanguage(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want II18nLanguage[youtube.I18nLanguage]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("$"),
					WithService(svc),
				},
			},
			want: &I18nLanguage{
				DefaultFields: &pkg.DefaultFields{
					Service:  &youtube.Service{},
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "$",
				},
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &I18nLanguage{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with empty string value",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &I18nLanguage{
				DefaultFields: &pkg.DefaultFields{},
				Hl:            "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewI18nLanguage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewI18nLanguage() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
