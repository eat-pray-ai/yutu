package videoAbuseReportReason

import (
	"reflect"
	"testing"
)

func TestNewVideoAbuseReportReason(t *testing.T) {
	type args struct {
		opt []Option
	}
	tests := []struct {
		name string
		args args
		want VideoAbuseReportReason
	}{
		{
			name: "TestNewVideoAbuseReportReason",
			args: args{
				opt: []Option{
					WithHL("hl"),
				},
			},
			want: &videoAbuseReportReason{
				hl: "hl",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoAbuseReportReason(tt.args.opt...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewVideoAbuseReportReason() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
