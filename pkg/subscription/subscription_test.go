package subscription

import (
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"reflect"
	"testing"
)

func TestNewSubscription(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Subscription
	}{
		{
			name: "TestNewSubscription",
			args: args{
				opts: []Option{
					WithIDs([]string{"id1", "id2"}),
					WithSubscriberChannelId("subscriberChannelId"),
					WithDescription("description"),
					WithChannelId("channelId"),
					WithForChannelId("forChannelId"),
					WithMaxResults(10),
					WithMine(utils.BoolPtr("true")),
					WithMyRecentSubscribers(utils.BoolPtr("false")),
					WithMySubscribers(utils.BoolPtr("false")),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOnBehalfOfContentOwnerChannel("contentOwnerChannel"),
					WithOrder("relevance"),
				},
			},
			want: &subscription{
				IDs:                           []string{"id1", "id2"},
				SubscriberChannelId:           "subscriberChannelId",
				Description:                   "description",
				ChannelId:                     "channelId",
				ForChannelId:                  "forChannelId",
				MaxResults:                    10,
				Mine:                          utils.BoolPtr("true"),
				MyRecentSubscribers:           utils.BoolPtr("false"),
				MySubscribers:                 utils.BoolPtr("false"),
				OnBehalfOfContentOwner:        "contentOwner",
				OnBehalfOfContentOwnerChannel: "contentOwnerChannel",
				Order:                         "relevance",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSubscription(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("NewSubscription() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
