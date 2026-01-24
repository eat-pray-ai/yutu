// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewChannel(t *testing.T) {
	type args struct {
		opts []Option
	}

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
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{
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
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(nil),
					WithMine(nil),
					WithMySubscribers(nil),
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(&managedByMeFalse),
					WithMine(&mineFalse),
					WithMySubscribers(&mySubscribersFalse),
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{
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
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-5),
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{
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
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{
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
					WithService(&youtube.Service{}),
				},
			},
			want: &Channel{
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
				if got := NewChannel(tt.args.opts...); !reflect.DeepEqual(got.(*Channel).CategoryId, tt.want.(*Channel).CategoryId) ||
					!reflect.DeepEqual(got.(*Channel).ForHandle, tt.want.(*Channel).ForHandle) ||
					!reflect.DeepEqual(got.(*Channel).ForUsername, tt.want.(*Channel).ForUsername) ||
					!reflect.DeepEqual(got.(*Channel).Hl, tt.want.(*Channel).Hl) ||
					!reflect.DeepEqual(got.(*Channel).Ids, tt.want.(*Channel).Ids) ||
					!reflect.DeepEqual(got.(*Channel).ManagedByMe, tt.want.(*Channel).ManagedByMe) ||
					!reflect.DeepEqual(got.(*Channel).MaxResults, tt.want.(*Channel).MaxResults) ||
					!reflect.DeepEqual(got.(*Channel).Mine, tt.want.(*Channel).Mine) ||
					!reflect.DeepEqual(got.(*Channel).MySubscribers, tt.want.(*Channel).MySubscribers) ||
					!reflect.DeepEqual(got.(*Channel).OnBehalfOfContentOwner, tt.want.(*Channel).OnBehalfOfContentOwner) ||
					!reflect.DeepEqual(got.(*Channel).Country, tt.want.(*Channel).Country) ||
					!reflect.DeepEqual(got.(*Channel).CustomUrl, tt.want.(*Channel).CustomUrl) ||
					!reflect.DeepEqual(got.(*Channel).DefaultLanguage, tt.want.(*Channel).DefaultLanguage) ||
					!reflect.DeepEqual(got.(*Channel).Description, tt.want.(*Channel).Description) ||
					!reflect.DeepEqual(got.(*Channel).Title, tt.want.(*Channel).Title) {
					t.Errorf("NewChannel() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
