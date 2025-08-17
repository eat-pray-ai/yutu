package activity

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/utils"
)

func TestNewActivity(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Activity
	}{
		{
			name: "TestNewActivity",
			args: args{
				opts: []Option{
					WithChannelId("ChannelId"),
					WithHome(utils.BoolPtr("true")),
					WithMaxResults(10),
					WithMine(utils.BoolPtr("true")),
					WithPublishedAfter("2021-01-01T00:00:00Z"),
					WithPublishedBefore("2021-01-31T00:00:00Z"),
					WithRegionCode("US"),
				},
			},
			want: &activity{
				ChannelId:       "ChannelId",
				Home:            utils.BoolPtr("true"),
				MaxResults:      10,
				Mine:            utils.BoolPtr("true"),
				PublishedAfter:  "2021-01-01T00:00:00Z",
				PublishedBefore: "2021-01-31T00:00:00Z",
				RegionCode:      "US",
			},
		},
		{
			name: "TestNewActivity",
			args: args{
				opts: []Option{
					WithHome(utils.BoolPtr("false")),
				},
			},
			want: &activity{
				Home: utils.BoolPtr("false"),
			},
		},
		{
			name: "TestNewActivity",
			args: args{
				opts: []Option{
					WithMaxResults(5),
				},
			},
			want: &activity{
				MaxResults: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewActivity(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewActivity() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
