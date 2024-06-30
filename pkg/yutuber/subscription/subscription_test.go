package subscription

import (
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
					WithID("id"),
					WithSubscriberChannelId("subscriberChannelId"),
					WithDescription("description"),
					WithChannelId("channelId"),
					WithForChannelId("forChannelId"),
					WithMaxResults(10),
					WithMine(true, true),
					WithMyRecentSubscribers(false, true),
					WithMySubscribers(false, true),
					WithOnBehalfOfContentOwner("contentOwner"),
					WithOnBehalfOfContentOwnerChannel("contentOwnerChannel"),
					WithOrder("relevance"),
				},
			},
			want: &subscription{
				ID:                            "id",
				SubscriberChannelId:           "subscriberChannelId",
				Description:                   "description",
				ChannelId:                     "channelId",
				ForChannelId:                  "forChannelId",
				MaxResults:                    10,
				Mine:                          &[]bool{true}[0],
				MyRecentSubscribers:           &[]bool{false}[0],
				MySubscribers:                 &[]bool{false}[0],
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
