// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package activity

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

func TestNewActivity(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
	homeTrue := true
	homeFalse := false
	mineTrue := true
	mineFalse := false

	tests := []struct {
		name string
		args args
		want IActivity[youtube.Activity]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithChannelId("test-channel-123"),
					WithHome(&homeTrue),
					WithMaxResults(50),
					WithMine(&mineTrue),
					WithPublishedAfter("2024-01-01T00:00:00Z"),
					WithPublishedBefore("2024-12-31T23:59:59Z"),
					WithRegionCode("US"),
					WithService(svc),
				},
			},
			want: &Activity{
				Fields:          &common.Fields{Service: svc},
				ChannelId:       "test-channel-123",
				Home:            &homeTrue,
				MaxResults:      50,
				Mine:            &mineTrue,
				PublishedAfter:  "2024-01-01T00:00:00Z",
				PublishedBefore: "2024-12-31T23:59:59Z",
				RegionCode:      "US",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Activity{Fields: &common.Fields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithHome(nil),
					WithMine(nil),
				},
			},
			want: &Activity{Fields: &common.Fields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithHome(&homeFalse),
					WithMine(&mineFalse),
				},
			},
			want: &Activity{
				Fields: &common.Fields{},
				Home:   &homeFalse,
				Mine:   &mineFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Activity{
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
			want: &Activity{
				Fields:     &common.Fields{},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithChannelId(""),
					WithPublishedAfter(""),
					WithPublishedBefore(""),
					WithRegionCode(""),
				},
			},
			want: &Activity{
				Fields:          &common.Fields{},
				ChannelId:       "",
				PublishedAfter:  "",
				PublishedBefore: "",
				RegionCode:      "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithChannelId("partial-channel"),
					WithMaxResults(25),
					WithRegionCode("UK"),
				},
			},
			want: &Activity{
				Fields:     &common.Fields{},
				ChannelId:  "partial-channel",
				MaxResults: 25,
				RegionCode: "UK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewActivity(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewActivity() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestActivity_Get(t *testing.T) {
	homeTrue := true
	mineTrue := true

	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get activities with channelId",
			opts: []Option{
				WithChannelId("channel-id"),
				WithMaxResults(2),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("channelId") != "channel-id" {
					t.Errorf(
						"expected channelId=channel-id, got %s",
						r.URL.Query().Get("channelId"),
					)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get activities with home",
			opts: []Option{
				WithHome(&homeTrue),
				WithMaxResults(2),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("home") != "true" {
					t.Errorf("expected home=true, got %s", r.URL.Query().Get("home"))
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get activities with mine",
			opts: []Option{
				WithMine(&mineTrue),
				WithMaxResults(2),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("mine") != "true" {
					t.Errorf("expected mine=true, got %s", r.URL.Query().Get("mine"))
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get activities with publishedAfter",
			opts: []Option{
				WithPublishedAfter("2024-01-01T00:00:00Z"),
				WithMaxResults(2),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("publishedAfter") != "2024-01-01T00:00:00Z" {
					t.Errorf(
						"expected publishedAfter=2024-01-01T00:00:00Z, got %s",
						r.URL.Query().Get("publishedAfter"),
					)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get activities with publishedBefore",
			opts: []Option{
				WithPublishedBefore("2024-12-31T23:59:59Z"),
				WithMaxResults(2),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("publishedBefore") != "2024-12-31T23:59:59Z" {
					t.Errorf(
						"expected publishedBefore=2024-12-31T23:59:59Z, got %s",
						r.URL.Query().Get("publishedBefore"),
					)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get activities with regionCode",
			opts: []Option{
				WithRegionCode("US"),
				WithMaxResults(2),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("regionCode") != "US" {
					t.Errorf(
						"expected regionCode=US, got %s", r.URL.Query().Get("regionCode"),
					)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get activities maxResults capping",
			opts: []Option{
				WithMaxResults(50),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("maxResults") != "20" {
					t.Errorf(
						"expected maxResults=20, got %s", r.URL.Query().Get("maxResults"),
					)
				}
			},
			wantLen: 2,
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
						{"id": "1", "snippet": {"title": "A1"}},
						{"id": "2", "snippet": {"title": "A2"}}
					],
					"nextPageToken": ""
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
				a := NewActivity(opts...)
				got, err := a.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Activity.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"Activity.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestActivity_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := range 20 {
				items[i] = `{"id": "id"}`
			}
			jsonItems := "[" + strings.Join(items, ",") + "]"
			_, _ = fmt.Fprintf(
				w, `{
                "items": %s,
                "nextPageToken": "page-2"
            }`, jsonItems,
			)
		} else if pageToken == "page-2" {
			_, _ = w.Write(
				[]byte(`{
                "items": [{"id": "21"}, {"id": "22"}],
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

	a := NewActivity(
		WithService(svc),
		WithChannelId("channel-id"),
		WithMaxResults(22),
	)
	got, err := a.Get()
	if err != nil {
		t.Errorf("Activity.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Activity.Get() got length = %v, want 22", len(got))
	}
}

func TestActivity_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "activity-1",
				"snippet": {
					"title": "Activity 1",
					"type": "upload",
					"publishedAt": "2024-01-01T00:00:00Z"
				}
			}
		],
		"nextPageToken": ""
	}`

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(mockResponse))
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
			name: "list activities json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithChannelId("channel-id"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list activities yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithChannelId("channel-id"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list activities table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithChannelId("channel-id"),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				a := NewActivity(tt.opts...)
				var buf bytes.Buffer
				if err := a.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Activity.List() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("Activity.List() output is empty")
				}
			},
		)
	}
}
