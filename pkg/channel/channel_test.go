// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewChannel(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
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
					WithParts([]string{"snippet", "contentDetails"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{
					Service:                svc,
					Parts:                  []string{"snippet", "contentDetails"},
					Output:                 "json",
					Ids:                    []string{"channel1", "channel2"},
					MaxResults:             100,
					Hl:                     "en",
					OnBehalfOfContentOwner: "owner123",
				},
				CategoryId:      "category123",
				ForHandle:       "@testhandle",
				ForUsername:     "testuser",
				ManagedByMe:     &managedByMeTrue,
				Mine:            &mineTrue,
				MySubscribers:   &mySubscribersTrue,
				Country:         "US",
				CustomUrl:       "testchannel",
				DefaultLanguage: "en",
				Description:     "Test channel description",
				Title:           "Test Channel",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{Service: svc},
			},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(nil),
					WithMine(nil),
					WithMySubscribers(nil),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{Service: svc},
			},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithChannelManagedByMe(&managedByMeFalse),
					WithMine(&mineFalse),
					WithMySubscribers(&mySubscribersFalse),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:        &common.Fields{Service: svc},
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
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{Service: svc, MaxResults: math.MaxInt64},
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-5),
					WithService(svc),
				},
			},
			want: &Channel{
				Fields: &common.Fields{Service: svc, MaxResults: 1},
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
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:          &common.Fields{Service: svc},
				CategoryId:      "",
				ForHandle:       "",
				ForUsername:     "",
				Country:         "",
				CustomUrl:       "",
				DefaultLanguage: "",
				Description:     "",
				Title:           "",
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
					WithService(svc),
				},
			},
			want: &Channel{
				Fields:  &common.Fields{Service: svc, Ids: []string{"channel1"}, MaxResults: 50},
				Title:   "My Channel",
				Country: "UK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewChannel(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewChannel() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestChannel_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get channels by categoryId",
			opts: []Option{
				WithCategoryId("category-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("categoryId") != "category-id" {
					t.Errorf(
						"expected categoryId=category-id, got %s",
						r.URL.Query().Get("categoryId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channels by forHandle",
			opts: []Option{
				WithForHandle("@handle"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("forHandle") != "@handle" {
					t.Errorf(
						"expected forHandle=@handle, got %s", r.URL.Query().Get("forHandle"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channels by forUsername",
			opts: []Option{
				WithForUsername("username"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("forUsername") != "username" {
					t.Errorf(
						"expected forUsername=username, got %s",
						r.URL.Query().Get("forUsername"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channels by id",
			opts: []Option{
				WithIds([]string{"id1", "id2"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				ids := r.URL.Query()["id"]
				if len(ids) == 1 && ids[0] == "id1,id2" {
					return
				}
				if len(ids) == 2 && ids[0] == "id1" && ids[1] == "id2" {
					return
				}
				t.Logf("ids: %v", ids)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channels managedByMe",
			opts: []Option{
				func(c *Channel) {
					b := true
					c.ManagedByMe = &b
				},
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("managedByMe") != "true" {
					t.Errorf(
						"expected managedByMe=true, got %s",
						r.URL.Query().Get("managedByMe"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channels by hl",
			opts: []Option{
				WithHl("en"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "en" {
					t.Errorf("expected hl=en, got %s", r.URL.Query().Get("hl"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get channels mine",
			opts: []Option{
				func(c *Channel) {
					b := true
					c.Mine = &b
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
			name: "get channels mySubscribers",
			opts: []Option{
				func(c *Channel) {
					b := true
					c.MySubscribers = &b
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
			name: "get channels with onBehalfOfContentOwner",
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
						{"id": "channel-1", "snippet": {"title": "Channel 1"}}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				c := NewChannel(opts...)
				got, err := c.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Channel.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"Channel.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestChannel_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("channel"))
	c := NewChannel(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := c.Get()
	if err != nil {
		t.Errorf("Channel.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Channel.Get() got length = %v, want 22", len(got))
	}
}

func TestChannel_List(t *testing.T) {
	common.RunListTest(
		t,
		`{"items": [{"id": "channel-1", "snippet": {"title": "Channel 1", "country": "US"}}]}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			c := NewChannel(
				WithService(svc),
				WithOutput(output),
				WithIds([]string{"channel-1"}),
			)
			return c.List
		},
	)
}

func TestChannel_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update channel",
			opts: []Option{
				WithIds([]string{"channel-id"}),
				WithTitle("New Title"),
				WithDescription("New Description"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "snippet" {
						t.Errorf("expected part=snippet, got %s", r.URL.Query().Get("part"))
					}
				} else if r.Method == "GET" {
					if r.URL.Query().Get("id") != "channel-id" {
						t.Errorf("expected id=channel-id, got %s", r.URL.Query().Get("id"))
					}
				} else {
					t.Errorf("unexpected method %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "update channel full",
			opts: []Option{
				WithIds([]string{"channel-id"}),
				WithCountry("US"),
				WithCustomUrl("my-url"),
				WithDefaultLanguage("en"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "snippet" {
						t.Errorf("expected part=snippet, got %s", r.URL.Query().Get("part"))
					}
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
							if r.Method == "GET" {
								_, _ = w.Write(
									[]byte(`{
						"items": [
							{"id": "channel-id", "snippet": {"title": "Old Title"}}
						]
					}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "channel-id", "snippet": {"title": "New Title"}}`))
							}
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				c := NewChannel(opts...)
				var buf bytes.Buffer
				if err := c.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Channel.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
