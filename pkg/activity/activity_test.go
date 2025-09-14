package activity

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewActivity(t *testing.T) {
	type args struct {
		opts []Option
	}

	homeTrue := true
	homeFalse := false
	mineTrue := true
	mineFalse := false

	tests := []struct {
		name string
		args args
		want Activity[youtube.Activity]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("test-channel-123"),
					WithHome(&homeTrue),
					WithMaxResults(50),
					WithMine(&mineTrue),
					WithPublishedAfter("2024-01-01T00:00:00Z"),
					WithPublishedBefore("2024-12-31T23:59:59Z"),
					WithRegionCode("US"),
					WithService(&youtube.Service{}),
				},
			},
			want: &activity{
				ChannelId:       "test-channel-123",
				Home:            &homeTrue,
				MaxResults:      50,
				Mine:            &mineTrue,
				PublishedAfter:  "2024-01-01T00:00:00Z",
				PublishedBefore: "2024-12-31T23:59:59Z",
				RegionCode:      "US",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &activity{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithHome(nil),
					WithMine(nil),
				},
			},
			want: &activity{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithHome(&homeFalse),
					WithMine(&mineFalse),
				},
			},
			want: &activity{
				Home: &homeFalse,
				Mine: &mineFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &activity{
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
			want: &activity{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithPublishedAfter(""),
					WithPublishedBefore(""),
					WithRegionCode(""),
				},
			},
			want: &activity{
				ChannelId:       "",
				PublishedAfter:  "",
				PublishedBefore: "",
				RegionCode:      "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("partial-channel"),
					WithMaxResults(25),
					WithRegionCode("UK"),
				},
			},
			want: &activity{
				ChannelId:  "partial-channel",
				MaxResults: 25,
				RegionCode: "UK",
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
