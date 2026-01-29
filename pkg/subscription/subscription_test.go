// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

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
				Fields: &common.Fields{
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
			want: &Subscription{Fields: &common.Fields{}},
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
			want: &Subscription{Fields: &common.Fields{}},
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
				Fields:              &common.Fields{},
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
			want: &Subscription{
				Fields:     &common.Fields{},
				MaxResults: 1,
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
				Fields:                        &common.Fields{},
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
				Fields:     &common.Fields{},
				ChannelId:  "myChannel",
				Title:      "My Subscription",
				MaxResults: 25,
				Order:      "alphabetical",
				Mine:       &mineTrue,
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

func TestSubscription_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get subscription by id",
			opts: []Option{
				WithIds([]string{"sub-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "sub-id" {
					t.Errorf("expected id=sub-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get subscription by channelId",
			opts: []Option{
				WithChannelId("channel-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("channelId") != "channel-id" {
					t.Errorf("expected channelId=channel-id, got %s", r.URL.Query().Get("channelId"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get subscription mine",
			opts: []Option{
				func(s *Subscription) {
					b := true
					s.Mine = &b
				},
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("mine") != "true" {
					t.Errorf("expected mine=true, got %s", r.URL.Query().Get("mine"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get subscription myRecentSubscribers",
			opts: []Option{
				func(s *Subscription) {
					b := true
					s.MyRecentSubscribers = &b
				},
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("myRecentSubscribers") != "true" {
					t.Errorf("expected myRecentSubscribers=true, got %s", r.URL.Query().Get("myRecentSubscribers"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get subscription mySubscribers",
			opts: []Option{
				func(s *Subscription) {
					b := true
					s.MySubscribers = &b
				},
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("mySubscribers") != "true" {
					t.Errorf("expected mySubscribers=true, got %s", r.URL.Query().Get("mySubscribers"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get subscription with onBehalfOfContentOwner",
			opts: []Option{
				WithOnBehalfOfContentOwner("owner-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
					t.Errorf("expected onBehalfOfContentOwner=owner-id, got %s", r.URL.Query().Get("onBehalfOfContentOwner"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{
					"items": [
						{"id": "sub-1", "snippet": {"title": "Subscription 1"}}
					]
				}`))
			}))
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
			s := NewSubscription(opts...)
			got, err := s.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscription.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Subscription.Get() got length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestSubscription_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := 0; i < 20; i++ {
				items[i] = fmt.Sprintf(`{"id": "sub-%d"}`, i)
			}
			w.Write([]byte(fmt.Sprintf(`{
				"items": [%s],
				"nextPageToken": "page-2"
			}`, strings.Join(items, ","))))
		} else if pageToken == "page-2" {
			w.Write([]byte(`{
				"items": [{"id": "sub-20"}, {"id": "sub-21"}],
				"nextPageToken": ""
			}`))
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

	s := NewSubscription(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := s.Get()
	if err != nil {
		t.Errorf("Subscription.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Subscription.Get() got length = %v, want 22", len(got))
	}
}

func TestSubscription_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"id": "sub-1",
					"snippet": {
						"resourceId": {
							"kind": "youtube#channel",
							"channelId": "channel-1"
						},
						"title": "Channel 1"
					}
				}
			]
		}`))
	}))
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
			name: "list subscriptions json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"sub-1"}),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list subscriptions yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"sub-1"}),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list subscriptions table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"sub-1"}),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSubscription(tt.opts...)
			var buf bytes.Buffer
			if err := s.List(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Subscription.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if buf.Len() == 0 {
				t.Errorf("Subscription.List() output is empty")
			}
		})
	}
}

func TestSubscription_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert subscription",
			opts: []Option{
				WithChannelId("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"id": "new-sub-id"}`))
			}))
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
			s := NewSubscription(opts...)
			var buf bytes.Buffer
			if err := s.Insert(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Subscription.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscription_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete subscription",
			opts: []Option{
				WithIds([]string{"sub-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
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
			s := NewSubscription(opts...)
			var buf bytes.Buffer
			if err := s.Delete(&buf); (err != nil) != tt.wantErr {
				t.Errorf("Subscription.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
