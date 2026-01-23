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
				opts: []Option{
					WithService(&youtube.Service{}),
				},
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
					WithService(&youtube.Service{}),
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
					WithService(&youtube.Service{}),
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
					WithService(&youtube.Service{}),
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
					WithService(&youtube.Service{}),
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
					WithService(&youtube.Service{}),
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
					WithService(&youtube.Service{}),
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
				if got := NewChannel(tt.args.opts...); !reflect.DeepEqual(got.(*channel).CategoryId, tt.want.(*channel).CategoryId) ||
					!reflect.DeepEqual(got.(*channel).ForHandle, tt.want.(*channel).ForHandle) ||
					!reflect.DeepEqual(got.(*channel).ForUsername, tt.want.(*channel).ForUsername) ||
					!reflect.DeepEqual(got.(*channel).Hl, tt.want.(*channel).Hl) ||
					!reflect.DeepEqual(got.(*channel).IDs, tt.want.(*channel).IDs) ||
					!reflect.DeepEqual(got.(*channel).ManagedByMe, tt.want.(*channel).ManagedByMe) ||
					!reflect.DeepEqual(got.(*channel).MaxResults, tt.want.(*channel).MaxResults) ||
					!reflect.DeepEqual(got.(*channel).Mine, tt.want.(*channel).Mine) ||
					!reflect.DeepEqual(got.(*channel).MySubscribers, tt.want.(*channel).MySubscribers) ||
					!reflect.DeepEqual(got.(*channel).OnBehalfOfContentOwner, tt.want.(*channel).OnBehalfOfContentOwner) ||
					!reflect.DeepEqual(got.(*channel).Country, tt.want.(*channel).Country) ||
					!reflect.DeepEqual(got.(*channel).CustomUrl, tt.want.(*channel).CustomUrl) ||
					!reflect.DeepEqual(got.(*channel).DefaultLanguage, tt.want.(*channel).DefaultLanguage) ||
					!reflect.DeepEqual(got.(*channel).Description, tt.want.(*channel).Description) ||
					!reflect.DeepEqual(got.(*channel).Title, tt.want.(*channel).Title) {
					t.Errorf("NewChannel() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
