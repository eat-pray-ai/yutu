// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewVideoAbuseReportReason(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IVideoAbuseReportReason[youtube.VideoAbuseReportReason]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHL("en"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[*].id"),
					WithService(svc),
				},
			},
			want: &VideoAbuseReportReason{
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
			want: &VideoAbuseReportReason{Fields: &common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHL(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &VideoAbuseReportReason{
				Fields: &common.Fields{
					Output: "", Jsonpath: "",
				},
				Hl: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHL("ja"),
					WithParts([]string{"snippet"}),
				},
			},
			want: &VideoAbuseReportReason{
				Fields: &common.Fields{Parts: []string{"snippet"}},
				Hl:     "ja",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoAbuseReportReason(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewVideoAbuseReportReason() = %v\nwant %v",
						tt.name, got, tt.want,
					)
				}
			},
		)
	}
}
