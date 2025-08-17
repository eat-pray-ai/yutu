package search

import (
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/utils"
)

func TestNewSearch(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Search
	}{
		{
			name: "TestNewSearch",
			args: args{
				opts: []Option{
					WithChannelId("channelId"),
					WithChannelType("show"),
					WithEventType("live"),
					WithForContentOwner(utils.BoolPtr("true")),
					WithForDeveloper(utils.BoolPtr("false")),
					WithForMine(utils.BoolPtr("false")),
					WithLocation("location"),
					WithLocationRadius("radius"),
					WithMaxResults(10),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOrder("rating"),
					WithTypes([]string{"video", "channel", "playlist"}),
				},
			},
			want: &search{
				ChannelId:              "channelId",
				ChannelType:            "show",
				EventType:              "live",
				ForContentOwner:        utils.BoolPtr("true"),
				ForDeveloper:           utils.BoolPtr("false"),
				ForMine:                utils.BoolPtr("false"),
				Location:               "location",
				LocationRadius:         "radius",
				MaxResults:             10,
				OnBehalfOfContentOwner: "contentOwner",
				Order:                  "rating",
				Types:                  []string{"video", "channel", "playlist"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSearch(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewSearch() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
