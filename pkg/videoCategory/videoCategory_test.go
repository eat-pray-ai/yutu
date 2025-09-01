package videoCategory

import (
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewVideoCategory(t *testing.T) {
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want VideoCategory
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"cat1", "cat2"}),
					WithHl("en"),
					WithRegionCode("US"),
					WithService(&youtube.Service{}),
				},
			},
			want: &videoCategory{
				IDs:        []string{"cat1", "cat2"},
				Hl:         "en",
				RegionCode: "US",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &videoCategory{},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
					WithRegionCode(""),
				},
			},
			want: &videoCategory{
				Hl:         "",
				RegionCode: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHl("ja"),
					WithRegionCode("JP"),
				},
			},
			want: &videoCategory{
				Hl:         "ja",
				RegionCode: "JP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoCategory(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewVideoCategory() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
