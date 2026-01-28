// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewChannel(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
	managedByMeTrue := true
	managedByMeFalse := false
	mineTrue := true
	mineFalse := false
	mySubscribersTrue := true
	mySubscribersFalse := false

	tests := []struct {
		name string
		args args
		want IChannel[youtube.Channel]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithCategoryId("category123"),
					WithForHandle("@testhandle"),
					WithForUsername("testuser"),
					WithHl("en"),
					WithIds([]string{"channel1", "channel2"}),
					WithChannelManagedByMe(&managedByMeTrue),
					WithMaxResults(100),
					WithMine(&mineTrue),
					WithMySubscribers(&mySubscribersTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithCountry("US"),
					WithCustomUrl("testchannel"),
					WithDefaultLanguage("en"),
					WithDescription("Test channel description"),
					WithTitle("Test Channel"),
					WithParts([]string{"snippet", "contentDetails"}),
					WithOutput("json"),
					WithJsonpath("$.items[0].id"),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"snippet", "contentDetails"},
					Output:   "json",
					Jsonpath: "$.items[0].id",
				},
				CategoryId:             "category123",
				ForHandle:              "@testhandle",
				ForUsername:            "testuser",
				Hl:                     "en",
				Ids:                    []string{"channel1", "channel2"},
				ManagedByMe:            &managedByMeTrue,
				MaxResults:             100,
				Mine:                   &mineTrue,
				MySubscribers:          &mySubscribersTrue,
				OnBehalfOfContentOwner: "owner123",
				Country:                "US",
				CustomUrl:              "testchannel",
				DefaultLanguage:        "en",
				Description:            "Test channel description",
				Title:                  "Test Channel",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{Service: svc},
			},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(nil),
					WithMine(nil),
					WithMySubscribers(nil),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{Service: svc},
			},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(&managedByMeFalse),
					WithMine(&mineFalse),
					WithMySubscribers(&mySubscribersFalse),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:        &common.Fields{Service: svc},
				ManagedByMe:   &managedByMeFalse,
				Mine:          &mineFalse,
				MySubscribers: &mySubscribersFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:     &common.Fields{Service: svc},
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-5),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:     &common.Fields{Service: svc},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithCategoryId(""),
					WithForHandle(""),
					WithForUsername(""),
					WithHl(""),
					WithOnBehalfOfContentOwner(""),
					WithCountry(""),
					WithCustomUrl(""),
					WithDefaultLanguage(""),
					WithDescription(""),
					WithTitle(""),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:                 &common.Fields{Service: svc},
				CategoryId:             "",
				ForHandle:              "",
				ForUsername:            "",
				Hl:                     "",
				OnBehalfOfContentOwner: "",
				Country:                "",
				CustomUrl:              "",
				DefaultLanguage:        "",
				Description:            "",
				Title:                  "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"channel1"}),
					WithTitle("My Channel"),
					WithCountry("UK"),
					WithMaxResults(50),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:     &common.Fields{Service: svc},
				Ids:        []string{"channel1"},
				Title:      "My Channel",
				Country:    "UK",
				MaxResults: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannel(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewChannel() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
