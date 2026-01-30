// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
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
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IVideo[youtube.Video]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"video1", "video2"}),
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
					WithParts([]string{"snippet", "contentDetails"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &Video{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"snippet", "contentDetails"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				Ids:                           []string{"video1", "video2"},
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
			want: &Video{Fields: &common.Fields{}},
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
			want: &Video{Fields: &common.Fields{}},
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
			want: &Video{
				Fields:              &common.Fields{},
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
			want: &Video{
				Fields:     &common.Fields{},
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
			want: &Video{
				Fields:     &common.Fields{},
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
			want: &Video{
				Fields:                        &common.Fields{},
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
			want: &Video{
				Fields:      &common.Fields{},
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
					t.Errorf("%s\nNewVideo() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestVideo_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get videos by id",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "video-id" {
					t.Errorf("expected id=video-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get videos myRating",
			opts: []Option{
				WithRating("like"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("myRating") != "like" {
					t.Errorf(
						"expected myRating=like, got %s", r.URL.Query().Get("myRating"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get videos chart",
			opts: []Option{
				WithChart("mostPopular"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("chart") != "mostPopular" {
					t.Errorf(
						"expected chart=mostPopular, got %s", r.URL.Query().Get("chart"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get videos with regionCode, maxHeight, maxWidth, onBehalfOfContentOwner",
			opts: []Option{
				WithRegionCode("US"),
				WithMaxHeight(1080),
				WithMaxWidth(1920),
				WithOnBehalfOfContentOwner("owner-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("regionCode") != "US" {
					t.Errorf(
						"expected regionCode=US, got %s", r.URL.Query().Get("regionCode"),
					)
				}
				if r.URL.Query().Get("maxHeight") != "1080" {
					t.Errorf(
						"expected maxHeight=1080, got %s", r.URL.Query().Get("maxHeight"),
					)
				}
				if r.URL.Query().Get("maxWidth") != "1920" {
					t.Errorf(
						"expected maxWidth=1920, got %s", r.URL.Query().Get("maxWidth"),
					)
				}
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write(
								[]byte(`{
					"items": [
						{"id": "video-1", "snippet": {"title": "Video 1"}}
					]
				}`),
							)
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				got, err := v.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Video.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("Video.Get() got length = %v, want %v", len(got), tt.wantLen)
				}
			},
		)
	}
}

func TestVideo_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := 0; i < 20; i++ {
				items[i] = fmt.Sprintf(`{"id": "video-%d"}`, i)
			}
			_, _ = w.Write(
				[]byte(fmt.Sprintf(
					`{
				"items": [%s],
				"nextPageToken": "page-2"
			}`, strings.Join(items, ","),
				)),
			)
		} else if pageToken == "page-2" {
			_, _ = w.Write(
				[]byte(`{
				"items": [{"id": "video-20"}, {"id": "video-21"}],
				"nextPageToken": ""
			}`),
			)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	svc, err := youtube.NewService(
		context.Background(),
		option.WithEndpoint(ts.URL),
		option.WithAPIKey("test-key"),
	)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	v := NewVideo(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := v.Get()
	if err != nil {
		t.Errorf("Video.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Video.Get() got length = %v, want 22", len(got))
	}
}

func TestVideo_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "video-1",
					"snippet": {
						"channelId": "channel-1",
						"title": "Video 1"
					},
					"statistics": {
						"viewCount": 100
					}
				}
			]
		}`),
				)
			},
		),
	)
	defer ts.Close()

	svc, err := youtube.NewService(
		context.Background(),
		option.WithEndpoint(ts.URL),
		option.WithAPIKey("test-key"),
	)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	tests := []struct {
		name    string
		opts    []Option
		output  string
		wantErr bool
	}{
		{
			name: "list videos json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"video-1"}),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list videos yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"video-1"}),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list videos table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"video-1"}),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				v := NewVideo(tt.opts...)
				var buf bytes.Buffer
				if err := v.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.List() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("Video.List() output is empty")
				}
			},
		)
	}
}

func TestVideo_Insert(t *testing.T) {
	autoLevelsTrue := true
	notifySubscribersTrue := true
	stabilizeTrue := true

	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert video",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithTitle("New Video"),
				WithPrivacy("public"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert video with options",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithTitle("New Video"),
				WithPrivacy("public"),
				WithAutoLevels(&autoLevelsTrue),
				WithNotifySubscribers(&notifySubscribersTrue),
				WithStabilize(&stabilizeTrue),
				WithOnBehalfOfContentOwner("owner-id"),
				WithOnBehalfOfContentOwnerChannel("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("autoLevels") != "true" {
					t.Errorf(
						"expected autoLevels=true, got %s", r.URL.Query().Get("autoLevels"),
					)
				}
				if r.URL.Query().Get("notifySubscribers") != "true" {
					t.Errorf(
						"expected notifySubscribers=true, got %s",
						r.URL.Query().Get("notifySubscribers"),
					)
				}
				if r.URL.Query().Get("stabilize") != "true" {
					t.Errorf(
						"expected stabilize=true, got %s", r.URL.Query().Get("stabilize"),
					)
				}
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
				if r.URL.Query().Get("onBehalfOfContentOwnerChannel") != "channel-id" {
					t.Errorf(
						"expected onBehalfOfContentOwnerChannel=channel-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwnerChannel"),
					)
				}
			},
			wantErr: false,
		},
	}

	err := os.WriteFile("test_video.mp4", []byte("dummy video content"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer func() {
		_ = os.Remove("test_video.mp4")
	}()

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"id": "new-video-id", "snippet": {"title": "New Video"}, "status": {"privacyStatus": "public"}}`))
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithTitle("Updated Title"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "snippet,status" {
						t.Errorf(
							"expected part=snippet,status, got %s", r.URL.Query().Get("part"),
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update video with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithTitle("Updated Title"),
				WithMaxResults(1),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
						t.Errorf(
							"expected onBehalfOfContentOwner=owner-id, got %s",
							r.URL.Query().Get("onBehalfOfContentOwner"),
						)
					}
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							if r.Method == "GET" {
								_, _ = w.Write(
									[]byte(`{
						"items": [
							{"id": "video-id", "snippet": {"title": "Old Title"}}
						]
					}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "video-id", "snippet": {"title": "Updated Title"}}`))
							}
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_Rate(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "rate video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithRating("like"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("id") != "video-id" {
					t.Errorf("expected id=video-id, got %s", r.URL.Query().Get("id"))
				}
				if r.URL.Query().Get("rating") != "like" {
					t.Errorf("expected rating=like, got %s", r.URL.Query().Get("rating"))
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.WriteHeader(http.StatusNoContent)
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Rate(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Rate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_GetRating(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "get rating",
			opts: []Option{
				WithIds([]string{"video-id"}),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "video-id" {
					t.Errorf("expected id=video-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantErr: false,
		},
		{
			name: "get rating with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write(
								[]byte(`{
					"items": [
						{"videoId": "video-id", "rating": "like"}
					]
				}`),
							)
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.GetRating(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.GetRating() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete video",
			opts: []Option{
				WithIds([]string{"video-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				if r.URL.Query().Get("id") != "video-id" {
					t.Errorf("expected id=video-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantErr: false,
		},
		{
			name: "delete video with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.WriteHeader(http.StatusNoContent)
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_ReportAbuse(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "report abuse",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithReasonId("reason-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "report abuse with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithReasonId("reason-id"),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf(
						"expected onBehalfOfContentOwner=owner-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.WriteHeader(http.StatusNoContent)
						},
					),
				)
				defer ts.Close()

				svc, err := youtube.NewService(
					context.Background(),
					option.WithEndpoint(ts.URL),
					option.WithAPIKey("test-key"),
				)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.ReportAbuse(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.ReportAbuse() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
