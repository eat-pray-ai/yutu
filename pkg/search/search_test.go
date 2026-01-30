// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package search

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
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
				Fields: &common.Fields{
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
			want: &Search{Fields: &common.Fields{}},
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
			want: &Search{Fields: &common.Fields{}},
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
				Fields:          &common.Fields{},
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
			want: &Search{
				Fields:     &common.Fields{},
				MaxResults: 1,
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
				Fields:                    &common.Fields{},
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
				Fields:     &common.Fields{},
				Q:          "golang tutorial",
				MaxResults: 25,
				Order:      "date",
				RegionCode: "UK",
				Types:      []string{"video"},
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

func TestSearch_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get search with query",
			opts: []Option{
				WithQ("test query"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("q") != "test query" {
					t.Errorf("expected q=test query, got %s", r.URL.Query().Get("q"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get search with channelId",
			opts: []Option{
				WithChannelId("channel-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("channelId") != "channel-id" {
					t.Errorf(
						"expected channelId=channel-id, got %s",
						r.URL.Query().Get("channelId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get search with types",
			opts: []Option{
				WithTypes([]string{"video", "playlist"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				gotType := r.URL.Query()["type"]
				joined := strings.Join(gotType, ",")
				if !strings.Contains(joined, "video") || !strings.Contains(
					joined, "playlist",
				) {
					t.Errorf(
						"expected type to contain video and playlist, got %v", gotType,
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get search with onBehalfOfContentOwner",
			opts: []Option{
				WithOnBehalfOfContentOwner("owner-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
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
						{"id": {"videoId": "video-1"}, "snippet": {"title": "Video 1"}}
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
				s := NewSearch(opts...)
				got, err := s.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Search.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("Search.Get() got length = %v, want %v", len(got), tt.wantLen)
				}
			},
		)
	}
}

func TestSearch_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := 0; i < 20; i++ {
				items[i] = fmt.Sprintf(`{"id": {"videoId": "video-%d"}}`, i)
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
				"items": [{"id": {"videoId": "video-20"}}, {"id": {"videoId": "video-21"}}],
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

	s := NewSearch(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := s.Get()
	if err != nil {
		t.Errorf("Search.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Search.Get() got length = %v, want 22", len(got))
	}
}

func TestSearch_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": {
						"kind": "youtube#video",
						"videoId": "video-1"
					},
					"snippet": {
						"title": "Video 1"
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
			name: "list search results json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithQ("test"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list search results yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithQ("test"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list search results table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithQ("test"),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := NewSearch(tt.opts...)
				var buf bytes.Buffer
				if err := s.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Search.List() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("Search.List() output is empty")
				}
			},
		)
	}
}
