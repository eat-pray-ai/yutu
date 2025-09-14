package video

import (
	"math"
	"reflect"
	"testing"

	"google.golang.org/api/youtube/v3"
)

func TestNewVideo(t *testing.T) {
	type args struct {
		opts []Option
	}

	autoLevelsTrue := true
	autoLevelsFalse := false
	forKidsTrue := true
	forKidsFalse := false
	embeddableTrue := true
	embeddableFalse := false
	stabilizeTrue := true
	stabilizeFalse := false
	notifySubscribersTrue := true
	notifySubscribersFalse := false
	publicStatsViewableTrue := true
	publicStatsViewableFalse := false

	tests := []struct {
		name string
		args args
		want Video[youtube.Video]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIDs([]string{"video1", "video2"}),
					WithAutoLevels(&autoLevelsTrue),
					WithFile("/path/to/video.mp4"),
					WithTitle("Test Video"),
					WithDescription("Test video description"),
					WithHl("en"),
					WithTags([]string{"tag1", "tag2", "tag3"}),
					WithLanguage("en"),
					WithLocale("en_US"),
					WithLicense("youtube"),
					WithThumbnail("/path/to/thumbnail.jpg"),
					WithRating("like"),
					WithChart("mostPopular"),
					WithChannelId("channel123"),
					WithComments("Test comments"),
					WithPlaylistId("playlist123"),
					WithCategory("22"),
					WithPrivacy("public"),
					WithForKids(&forKidsTrue),
					WithEmbeddable(&embeddableTrue),
					WithPublishAt("2024-12-31T23:59:59Z"),
					WithRegionCode("US"),
					WithReasonId("reason123"),
					WithSecondaryReasonId("secondaryReason123"),
					WithStabilize(&stabilizeTrue),
					WithMaxHeight(1080),
					WithMaxWidth(1920),
					WithMaxResults(50),
					WithNotifySubscribers(&notifySubscribersTrue),
					WithPublicStatsViewable(&publicStatsViewableTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithService(&youtube.Service{}),
				},
			},
			want: &video{
				IDs:                           []string{"video1", "video2"},
				AutoLevels:                    &autoLevelsTrue,
				File:                          "/path/to/video.mp4",
				Title:                         "Test Video",
				Description:                   "Test video description",
				Hl:                            "en",
				Tags:                          []string{"tag1", "tag2", "tag3"},
				Language:                      "en",
				Locale:                        "en_US",
				License:                       "youtube",
				Thumbnail:                     "/path/to/thumbnail.jpg",
				Rating:                        "like",
				Chart:                         "mostPopular",
				ChannelId:                     "channel123",
				Comments:                      "Test comments",
				PlaylistId:                    "playlist123",
				CategoryId:                    "22",
				Privacy:                       "public",
				ForKids:                       &forKidsTrue,
				Embeddable:                    &embeddableTrue,
				PublishAt:                     "2024-12-31T23:59:59Z",
				RegionCode:                    "US",
				ReasonId:                      "reason123",
				SecondaryReasonId:             "secondaryReason123",
				Stabilize:                     &stabilizeTrue,
				MaxHeight:                     1080,
				MaxWidth:                      1920,
				MaxResults:                    50,
				NotifySubscribers:             &notifySubscribersTrue,
				PublicStatsViewable:           &publicStatsViewableTrue,
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &video{},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithAutoLevels(nil),
					WithForKids(nil),
					WithEmbeddable(nil),
					WithStabilize(nil),
					WithNotifySubscribers(nil),
					WithPublicStatsViewable(nil),
				},
			},
			want: &video{},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithAutoLevels(&autoLevelsFalse),
					WithForKids(&forKidsFalse),
					WithEmbeddable(&embeddableFalse),
					WithStabilize(&stabilizeFalse),
					WithNotifySubscribers(&notifySubscribersFalse),
					WithPublicStatsViewable(&publicStatsViewableFalse),
				},
			},
			want: &video{
				AutoLevels:          &autoLevelsFalse,
				ForKids:             &forKidsFalse,
				Embeddable:          &embeddableFalse,
				Stabilize:           &stabilizeFalse,
				NotifySubscribers:   &notifySubscribersFalse,
				PublicStatsViewable: &publicStatsViewableFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &video{
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &video{
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithFile(""),
					WithTitle(""),
					WithDescription(""),
					WithHl(""),
					WithLanguage(""),
					WithLocale(""),
					WithLicense(""),
					WithThumbnail(""),
					WithRating(""),
					WithChart(""),
					WithChannelId(""),
					WithComments(""),
					WithPlaylistId(""),
					WithCategory(""),
					WithPrivacy(""),
					WithPublishAt(""),
					WithRegionCode(""),
					WithReasonId(""),
					WithSecondaryReasonId(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &video{
				File:                          "",
				Title:                         "",
				Description:                   "",
				Hl:                            "",
				Language:                      "",
				Locale:                        "",
				License:                       "",
				Thumbnail:                     "",
				Rating:                        "",
				Chart:                         "",
				ChannelId:                     "",
				Comments:                      "",
				PlaylistId:                    "",
				CategoryId:                    "",
				Privacy:                       "",
				PublishAt:                     "",
				RegionCode:                    "",
				ReasonId:                      "",
				SecondaryReasonId:             "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Video"),
					WithDescription("A great video"),
					WithTags([]string{"tutorial", "golang"}),
					WithPrivacy("private"),
					WithMaxResults(25),
					WithForKids(&forKidsFalse),
				},
			},
			want: &video{
				Title:       "My Video",
				Description: "A great video",
				Tags:        []string{"tutorial", "golang"},
				Privacy:     "private",
				MaxResults:  25,
				ForKids:     &forKidsFalse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideo(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewVideo() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
