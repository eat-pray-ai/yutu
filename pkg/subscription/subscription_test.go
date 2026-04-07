// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package subscription

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
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
					WithService(svc),
				},
			},
			want: &Subscription{
				Fields: &common.Fields{
					Service:                svc,
					Parts:                  []string{"snippet", "contentDetails"},
					Output:                 "json",
					Ids:                    []string{"sub1", "sub2"},
					MaxResults:             50,
					ChannelId:              "channel123",
					OnBehalfOfContentOwner: "owner123",
				},
				SubscriberChannelId:           "subscriber123",
				Description:                   "Test subscription description",
				ForChannelId:                  "forChannel123",
				Mine:                          &mineTrue,
				MyRecentSubscribers:           &myRecentSubscribersTrue,
				MySubscribers:                 &mySubscribersTrue,
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
				Fields: &common.Fields{MaxResults: math.MaxInt64},
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
				Fields: &common.Fields{MaxResults: 1},
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
				Fields: &common.Fields{
					ChannelId:              "",
					OnBehalfOfContentOwner: "",
				},
				SubscriberChannelId:           "",
				Description:                   "",
				ForChannelId:                  "",
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
				Fields: &common.Fields{
					ChannelId:  "myChannel",
					MaxResults: 25,
				},
				Title: "My Subscription",
				Order: "alphabetical",
				Mine:  &mineTrue,
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
					t.Errorf(
						"expected myRecentSubscribers=true, got %s",
						r.URL.Query().Get("myRecentSubscribers"),
					)
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
					t.Errorf(
						"expected mySubscribers=true, got %s",
						r.URL.Query().Get("mySubscribers"),
					)
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
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write(
								[]byte(`{
					"items": [
						{"id": "sub-1", "snippet": {"title": "Subscription 1"}}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewSubscription(opts...)
				got, err := s.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Subscription.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"Subscription.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestSubscription_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("sub"))
	s := NewSubscription(WithService(svc), WithMaxResults(22))
	got, err := s.Get()
	if err != nil {
		t.Errorf("Subscription.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Subscription.Get() got length = %v, want 22", len(got))
	}
}

func TestSubscription_List(t *testing.T) {
	common.RunListTest(t, `{
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
		}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			s := NewSubscription(WithService(svc), WithOutput(output), WithIds([]string{"sub-1"}))
			return s.List
		},
	)
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
				WithTitle("my-title"),
				WithDescription("my-description"),
				WithSubscriberChannelId("subscriber-channel-id"),
			},
			verify: func(r *http.Request) {
				// ensure HTTP method
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				// decode request body and verify payload
				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						ResourceId struct {
							ChannelId string `json:"channelId"`
						} `json:"resourceId"`
						Title       string `json:"title"`
						Description string `json:"description"`
						ChannelId   string `json:"channelId"`
					} `json:"snippet"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Snippet.ResourceId.ChannelId != "channel-id" {
					t.Errorf(
						"expected snippet.resourceId.channelId=channel-id, got %s",
						body.Snippet.ResourceId.ChannelId,
					)
				}
				if body.Snippet.Title != "my-title" {
					t.Errorf("expected snippet.title=my-title, got %s", body.Snippet.Title)
				}
				if body.Snippet.Description != "my-description" {
					t.Errorf(
						"expected snippet.description=my-description, got %s",
						body.Snippet.Description,
					)
				}
				if body.Snippet.ChannelId != "subscriber-channel-id" {
					t.Errorf(
						"expected snippet.channelId=subscriber-channel-id, got %s",
						body.Snippet.ChannelId,
					)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"id": "new-sub-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewSubscription(opts...)
				var buf bytes.Buffer
				if err := s.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"Subscription.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
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
		t.Run(
			tt.name, func(t *testing.T) {
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.WriteHeader(http.StatusNoContent)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewSubscription(opts...)
				var buf bytes.Buffer
				if err := s.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"Subscription.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
