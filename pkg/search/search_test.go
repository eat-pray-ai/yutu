// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"math"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"google.golang.org/api/youtube/v3"
)

func TestNewSearch(t *testing.T) {
	type args struct {
		opts []Option
	}

	forContentOwnerTrue := true
	forContentOwnerFalse := false
	forDeveloperTrue := true
	forDeveloperFalse := false
	forMineTrue := true
	forMineFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ISearch[youtube.SearchResult]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("channel123"),
					WithChannelType("any"),
					WithEventType("live"),
					WithForContentOwner(&forContentOwnerTrue),
					WithForDeveloper(&forDeveloperTrue),
					WithForMine(&forMineTrue),
					WithLocation("37.42307,-122.08427"),
					WithLocationRadius("50km"),
					WithMaxResults(50),
					WithOnBehalfOfContentOwner("owner123"),
					WithOrder("relevance"),
					WithPublishedAfter("2024-01-01T00:00:00Z"),
					WithPublishedBefore("2024-12-31T23:59:59Z"),
					WithQ("test search query"),
					WithRegionCode("US"),
					WithRelevanceLanguage("en"),
					WithSafeSearch("moderate"),
					WithTopicId("/m/04rlf"),
					WithTypes([]string{"video", "channel", "playlist"}),
					WithVideoCaption("closedCaption"),
					WithVideoCategoryId("10"),
					WithVideoDefinition("high"),
					WithVideoDimension("2d"),
					WithVideoDuration("medium"),
					WithVideoEmbeddable("true"),
					WithVideoLicense("youtube"),
					WithVideoPaidProductPlacement("true"),
					WithVideoSyndicated("true"),
					WithVideoType("movie"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &Search{
				DefaultFields: &pkg.DefaultFields{
					Service:  svc,
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				ChannelId:                 "channel123",
				ChannelType:               "any",
				EventType:                 "live",
				ForContentOwner:           &forContentOwnerTrue,
				ForDeveloper:              &forDeveloperTrue,
				ForMine:                   &forMineTrue,
				Location:                  "37.42307,-122.08427",
				LocationRadius:            "50km",
				MaxResults:                50,
				OnBehalfOfContentOwner:    "owner123",
				Order:                     "relevance",
				PublishedAfter:            "2024-01-01T00:00:00Z",
				PublishedBefore:           "2024-12-31T23:59:59Z",
				Q:                         "test search query",
				RegionCode:                "US",
				RelevanceLanguage:         "en",
				SafeSearch:                "moderate",
				TopicId:                   "/m/04rlf",
				Types:                     []string{"video", "channel", "playlist"},
				VideoCaption:              "closedCaption",
				VideoCategoryId:           "10",
				VideoDefinition:           "high",
				VideoDimension:            "2d",
				VideoDuration:             "medium",
				VideoEmbeddable:           "true",
				VideoLicense:              "youtube",
				VideoPaidProductPlacement: "true",
				VideoSyndicated:           "true",
				VideoType:                 "movie",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Search{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithForContentOwner(nil),
					WithForDeveloper(nil),
					WithForMine(nil),
				},
			},
			want: &Search{DefaultFields: &pkg.DefaultFields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithForContentOwner(&forContentOwnerFalse),
					WithForDeveloper(&forDeveloperFalse),
					WithForMine(&forMineFalse),
				},
			},
			want: &Search{
				DefaultFields:   &pkg.DefaultFields{},
				ForContentOwner: &forContentOwnerFalse,
				ForDeveloper:    &forDeveloperFalse,
				ForMine:         &forMineFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Search{
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
			want: &Search{
				DefaultFields: &pkg.DefaultFields{},
				MaxResults:    1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithChannelType(""),
					WithEventType(""),
					WithLocation(""),
					WithLocationRadius(""),
					WithOnBehalfOfContentOwner(""),
					WithOrder(""),
					WithPublishedAfter(""),
					WithPublishedBefore(""),
					WithQ(""),
					WithRegionCode(""),
					WithRelevanceLanguage(""),
					WithSafeSearch(""),
					WithTopicId(""),
					WithVideoCaption(""),
					WithVideoCategoryId(""),
					WithVideoDefinition(""),
					WithVideoDimension(""),
					WithVideoDuration(""),
					WithVideoEmbeddable(""),
					WithVideoLicense(""),
					WithVideoPaidProductPlacement(""),
					WithVideoSyndicated(""),
					WithVideoType(""),
				},
			},
			want: &Search{
				DefaultFields:             &pkg.DefaultFields{},
				ChannelId:                 "",
				ChannelType:               "",
				EventType:                 "",
				Location:                  "",
				LocationRadius:            "",
				OnBehalfOfContentOwner:    "",
				Order:                     "",
				PublishedAfter:            "",
				PublishedBefore:           "",
				Q:                         "",
				RegionCode:                "",
				RelevanceLanguage:         "",
				SafeSearch:                "",
				TopicId:                   "",
				VideoCaption:              "",
				VideoCategoryId:           "",
				VideoDefinition:           "",
				VideoDimension:            "",
				VideoDuration:             "",
				VideoEmbeddable:           "",
				VideoLicense:              "",
				VideoPaidProductPlacement: "",
				VideoSyndicated:           "",
				VideoType:                 "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithQ("golang tutorial"),
					WithMaxResults(25),
					WithOrder("date"),
					WithRegionCode("UK"),
					WithTypes([]string{"video"}),
				},
			},
			want: &Search{
				DefaultFields: &pkg.DefaultFields{},
				Q:             "golang tutorial",
				MaxResults:    25,
				Order:         "date",
				RegionCode:    "UK",
				Types:         []string{"video"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSearch(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewSearch() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}
