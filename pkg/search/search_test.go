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
					WithForContentOwner(true, true),
					WithForDeveloper(false, true),
					WithForMine(false, true),
					WithLocation("location"),
					WithLocationRadius("radius"),
					WithMaxResults(10),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOrder("rating"),
				},
			},
			want: &search{
				ChannelId:              "channelId",
				ChannelType:            "show",
				EventType:              "live",
				ForContentOwner:        &[]bool{true}[0],
				ForDeveloper:           &[]bool{false}[0],
				ForMine:                &[]bool{false}[0],
				Location:               "location",
				LocationRadius:         "radius",
				MaxResults:             10,
				OnBehalfOfContentOwner: "contentOwner",
				Order:                  "rating",
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
