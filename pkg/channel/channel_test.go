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
		want Channel[youtube.Channel]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithCategoryId("category123"),
					WithForHandle("@testhandle"),
					WithForUsername("testuser"),
					WithHl("en"),
					WithIDs([]string{"channel1", "channel2"}),
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
			want: &channel{
				CategoryId:             "category123",
				ForHandle:              "@testhandle",
				ForUsername:            "testuser",
				Hl:                     "en",
				IDs:                    []string{"channel1", "channel2"},
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
				opts: []Option{},
			},
			want: &channel{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(nil),
					WithMine(nil),
					WithMySubscribers(nil),
				},
			},
			want: &channel{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(&managedByMeFalse),
					WithMine(&mineFalse),
					WithMySubscribers(&mySubscribersFalse),
				},
			},
			want: &channel{
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
				},
			},
			want: &channel{
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-5),
				},
			},
			want: &channel{
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
				},
			},
			want: &channel{
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
					WithIDs([]string{"channel1"}),
					WithTitle("My Channel"),
					WithCountry("UK"),
					WithMaxResults(50),
				},
			},
			want: &channel{
				IDs:        []string{"channel1"},
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
					t.Errorf("NewChannel() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
