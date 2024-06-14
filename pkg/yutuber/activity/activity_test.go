package activity

import (
	"reflect"
	"testing"
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
					WithChannelId("channelId"),
					WithHome(true, true),
					WithMaxResults(10),
					WithMine(true, true),
					WithPublishedAfter("2021-01-01T00:00:00Z"),
					WithPublishedBefore("2021-01-31T00:00:00Z"),
					WithRegionCode("US"),
				},
			},
			want: &activity{
				channelId:       "channelId",
				home:            &[]bool{true}[0],
				maxResults:      10,
				mine:            &[]bool{true}[0],
				publishedAfter:  "2021-01-01T00:00:00Z",
				publishedBefore: "2021-01-31T00:00:00Z",
				regionCode:      "US",
			},
		},
		{
			name: "TestNewActivity",
			args: args{
				opts: []Option{
					WithHome(false, true),
				},
			},
			want: &activity{
				home: &[]bool{false}[0],
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
				maxResults: 5,
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
