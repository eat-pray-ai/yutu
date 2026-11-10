// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package channel

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

func TestNewChannel(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

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
					WithFor("managedByMe"),
					WithMaxResults(100),
					WithOnBehalfOfContentOwner("owner123"),
					WithCountry(new("US")),
					WithDefaultLanguage(new("en")),
					WithDescription(new("Test channel description")),
					WithTitle("Test Channel"),
					WithParts([]string{"snippet", "contentDetails"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &Channel{
				Service:                svc,
				Parts:                  []string{"snippet", "contentDetails"},
				Output:                 "json",
				Ids:                    []string{"channel1", "channel2"},
				MaxResults:             100,
				Hl:                     "en",
				OnBehalfOfContentOwner: "owner123",
				CategoryId:             "category123",
				ForHandle:              "@testhandle",
				ForUsername:            "testuser",
				For:                    "managedByMe",
				Country:                new("US"),
				DefaultLanguage:        new("en"),
				Description:            new("Test channel description"),
				Title:                  "Test Channel",
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
				Service: svc,
			},
		},
		{
			name: "with empty for option",
			args: args{
				opts: []Option{
					WithFor(""),
					WithService(svc),
				},
			},
			want: &Channel{
				Service: svc,
			},
		},
		{
			name: "with for mine",
			args: args{
				opts: []Option{
					WithFor("mine"),
					WithService(svc),
				},
			},
			want: &Channel{
				Service: svc,
				For:     "mine",
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
				Service: svc, MaxResults: math.MaxInt64,
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
				Service: svc, MaxResults: 1,
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
					WithCountry(new("")),
					WithDefaultLanguage(new("")),
					WithDescription(new("")),
					WithTitle(""),
					WithService(svc),
				},
			},
			want: &Channel{
				Service:         svc,
				CategoryId:      "",
				ForHandle:       "",
				ForUsername:     "",
				Country:         new(""),
				DefaultLanguage: new(""),
				Description:     new(""),
				Title:           "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"channel1"}),
					WithTitle("My Channel"),
					WithCountry(new("UK")),
					WithMaxResults(50),
					WithService(svc),
				},
			},
			want: &Channel{
				Service: svc, Ids: []string{"channel1"}, MaxResults: 50,
				Title:   "My Channel",
				Country: new("UK"),
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
				WithFor("managedByMe"),
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
				WithFor("mine"),
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
				WithFor("mySubscribers"),
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
				WithDescription(new("New Description")),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				switch r.Method {
				case "PUT":
					if r.URL.Query().Get("part") != "brandingSettings" {
						t.Errorf("expected part=brandingSettings, got %s", r.URL.Query().Get("part"))
					}
				case "GET":
					if r.URL.Query().Get("id") != "channel-id" {
						t.Errorf("expected id=channel-id, got %s", r.URL.Query().Get("id"))
					}
				default:
					t.Errorf("unexpected method %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			// The CLI and MCP tool never pass WithMaxResults; the lookup must still find the channel.
			name: "update channel without max results",
			opts: []Option{
				WithIds([]string{"channel-id"}),
				WithDescription(new("New Description")),
			},
			wantErr: false,
		},
		{
			name: "update channel full",
			opts: []Option{
				WithIds([]string{"channel-id"}),
				WithCountry(new("US")),
				WithDefaultLanguage(new("en")),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "brandingSettings" {
						t.Errorf("expected part=brandingSettings, got %s", r.URL.Query().Get("part"))
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

// YouTube only accepts brandingSettings (not snippet) on channels.update, and the
// update must keep branding fields the caller did not touch.
func TestChannel_Update_BrandingSettingsBody(t *testing.T) {
	var body map[string]any
	svc := common.NewTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "GET" {
			if got := r.URL.Query().Get("part"); got != "brandingSettings" {
				t.Errorf("GET part = %q, want brandingSettings", got)
			}
			_, _ = w.Write([]byte(`{"items": [{"id": "channel-id", "brandingSettings": {"channel": {"title": "Old", "keywords": "ai agents"}}}]}`))
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode PUT body: %v", err)
		}
		_, _ = w.Write([]byte(`{"id": "channel-id"}`))
	}))
	c := NewChannel(
		WithService(svc),
		WithIds([]string{"channel-id"}),
		WithDescription(new("New Description")),
		WithCountry(new("DE")),
		WithDefaultLanguage(new("en")),
	)
	if err := c.Update(&bytes.Buffer{}); err != nil {
		t.Fatalf("Channel.Update() error = %v", err)
	}
	ch := body["brandingSettings"].(map[string]any)["channel"].(map[string]any)
	for k, want := range map[string]string{"description": "New Description", "country": "DE", "defaultLanguage": "en", "keywords": "ai agents", "title": "Old"} {
		if ch[k] != want {
			t.Errorf("brandingSettings.channel.%s = %v, want %q", k, ch[k], want)
		}
	}
	if _, ok := body["snippet"]; ok {
		t.Errorf("PUT body must not send snippet")
	}
}
