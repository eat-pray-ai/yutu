// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewSubscription(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false
	myRecentSubscribersTrue := true
	myRecentSubscribersFalse := false
	mySubscribersTrue := true
	mySubscribersFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ISubscription[youtube.Subscription]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"sub1", "sub2"}),
					WithSubscriberChannelId("subscriber123"),
					WithDescription("Test subscription description"),
					WithChannelId("channel123"),
					WithForChannelId("forChannel123"),
					WithMaxResults(50),
					WithMine(&mineTrue),
					WithMyRecentSubscribers(&myRecentSubscribersTrue),
					WithMySubscribers(&mySubscribersTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithOrder("relevance"),
					WithTitle("Test Subscription"),
					WithParts([]string{"snippet", "contentDetails"}),
					WithOutput("json"),
					WithJsonpath("$.items[0].id"),
					WithService(svc),
				},
			},
			want: &Subscription{
				DefaultFields: &pkg.DefaultFields{
					Service:  svc,
					Parts:    []string{"snippet", "contentDetails"},
					Output:   "json",
					Jsonpath: "$.items[0].id",
				},
				Ids:                           []string{"sub1", "sub2"},
				SubscriberChannelId:           "subscriber123",
				Description:                   "Test subscription description",
				ChannelId:                     "channel123",
				ForChannelId:                  "forChannel123",
				MaxResults:                    50,
				Mine:                          &mineTrue,
				MyRecentSubscribers:           &myRecentSubscribersTrue,
				MySubscribers:                 &mySubscribersTrue,
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
				Order:                         "relevance",
				Title:                         "Test Subscription",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Subscription{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
					WithMyRecentSubscribers(nil),
					WithMySubscribers(nil),
				},
			},
			want: &Subscription{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
					WithMyRecentSubscribers(&myRecentSubscribersFalse),
					WithMySubscribers(&mySubscribersFalse),
				},
			},
			want: &Subscription{
				DefaultFields:       &pkg.DefaultFields{},
				Mine:                &mineFalse,
				MyRecentSubscribers: &myRecentSubscribersFalse,
				MySubscribers:       &mySubscribersFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Subscription{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &Subscription{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithSubscriberChannelId(""),
					WithDescription(""),
					WithChannelId(""),
					WithForChannelId(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
					WithOrder(""),
					WithTitle(""),
				},
			},
			want: &Subscription{
				DefaultFields:                 &pkg.DefaultFields{},
				SubscriberChannelId:           "",
				Description:                   "",
				ChannelId:                     "",
				ForChannelId:                  "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
				Order:                         "",
				Title:                         "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("myChannel"),
					WithTitle("My Subscription"),
					WithMaxResults(25),
					WithOrder("alphabetical"),
					WithMine(&mineTrue),
				},
			},
			want: &Subscription{
				DefaultFields: &pkg.DefaultFields{},
				ChannelId:     "myChannel",
				Title:         "My Subscription",
				MaxResults:    25,
				Order:         "alphabetical",
				Mine:          &mineTrue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSubscription(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewSubscription() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
