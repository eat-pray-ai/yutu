// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewVideoAbuseReportReason(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want VideoAbuseReportReason[youtube.VideoAbuseReportReason]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHL("en"),
					WithService(&youtube.Service{}),
				},
			},
			want: &videoAbuseReportReason{
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &videoAbuseReportReason{},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHL(""),
				},
			},
			want: &videoAbuseReportReason{
				Hl: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHL("ja"),
				},
			},
			want: &videoAbuseReportReason{
				Hl: "ja",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoAbuseReportReason(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewVideoAbuseReportReason() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
