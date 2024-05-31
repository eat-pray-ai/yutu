package search

import (
	"reflect"
	"testing"
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
					WithForContentOwner("contentOwner"),
					WithForDeveloper("developer"),
					WithForMine("true"),
					WithLocation("location"),
					WithLocationRadius("radius"),
					WithMaxResults(10),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOrder("rating"),
				},
			},
			want: &search{
				channelId:              "channelId",
				channelType:            "show",
				eventType:              "live",
				forContentOwner:        "contentOwner",
				forDeveloper:           "developer",
				forMine:                "true",
				location:               "location",
				locationRadius:         "radius",
				maxResults:             10,
				onBehalfOfContentOwner: "contentOwner",
				order:                  "rating",
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
