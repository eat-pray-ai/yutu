// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlist

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestNewPlaylist(t *testing.T) {
	type args struct {
		opts []Option
	}

	mineTrue := true
	mineFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IPlaylist[youtube.Playlist]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"playlist1", "playlist2"}),
					WithTitle("Test Playlist"),
					WithDescription("Test playlist description"),
					WithTags([]string{"tag1", "tag2", "tag3"}),
					WithLanguage("en"),
					WithChannelId("channel123"),
					WithPrivacy("public"),
					WithHl("en"),
					WithMaxResults(50),
					WithMine(&mineTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[0].snippet.title"),
					WithService(svc),
				},
			},
			want: &Playlist{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[0].snippet.title",
				},
				Ids:                           []string{"playlist1", "playlist2"},
				Title:                         "Test Playlist",
				Description:                   "Test playlist description",
				Tags:                          []string{"tag1", "tag2", "tag3"},
				Language:                      "en",
				ChannelId:                     "channel123",
				Privacy:                       "public",
				Hl:                            "en",
				MaxResults:                    50,
				Mine:                          &mineTrue,
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Playlist{Fields: &common.Fields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithMine(nil),
				},
			},
			want: &Playlist{Fields: &common.Fields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithMine(&mineFalse),
				},
			},
			want: &Playlist{
				Fields: &common.Fields{},
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
			want: &Playlist{
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
			want: &Playlist{
				Fields:     &common.Fields{},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithTitle(""),
					WithDescription(""),
					WithLanguage(""),
					WithChannelId(""),
					WithPrivacy(""),
					WithHl(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &Playlist{
				Fields:                        &common.Fields{},
				Title:                         "",
				Description:                   "",
				Language:                      "",
				ChannelId:                     "",
				Privacy:                       "",
				Hl:                            "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Playlist"),
					WithDescription("A great playlist"),
					WithPrivacy("private"),
					WithMaxResults(25),
				},
			},
			want: &Playlist{
				Fields:      &common.Fields{},
				Title:       "My Playlist",
				Description: "A great playlist",
				Privacy:     "private",
				MaxResults:  25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylist(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewPlaylist() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestPlaylist_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get playlists by id",
			opts: []Option{
				WithIds([]string{"playlist-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "playlist-id" {
					t.Errorf("expected id=playlist-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get playlists by channelId",
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
			name: "get playlists mine",
			opts: []Option{
				func(p *Playlist) {
					b := true
					p.Mine = &b
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
			name: "get playlists with onBehalfOfContentOwner",
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
		{
			name: "get playlists with hl",
			opts: []Option{
				WithHl("es"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "es" {
					t.Errorf("expected hl=es, got %s", r.URL.Query().Get("hl"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get playlists with onBehalfOfContentOwnerChannel",
			opts: []Option{
				WithOnBehalfOfContentOwnerChannel("channel-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwnerChannel") != "channel-id" {
					t.Errorf(
						"expected onBehalfOfContentOwnerChannel=channel-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwnerChannel"),
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
						{"id": "playlist-1", "snippet": {"title": "Playlist 1"}}
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
				p := NewPlaylist(opts...)
				got, err := p.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Playlist.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"Playlist.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestPlaylist_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := range 20 {
				items[i] = fmt.Sprintf(`{"id": "playlist-%d"}`, i)
			}
			_, _ = w.Write(
				fmt.Appendf(nil,
					`{
				"items": [%s],
				"nextPageToken": "page-2"
			}`, strings.Join(items, ","),
				),
			)
		} else if pageToken == "page-2" {
			_, _ = w.Write(
				[]byte(`{
				"items": [{"id": "playlist-20"}, {"id": "playlist-21"}],
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

	p := NewPlaylist(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := p.Get()
	if err != nil {
		t.Errorf("Playlist.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Playlist.Get() got length = %v, want 22", len(got))
	}
}

func TestPlaylist_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "playlist-1",
					"snippet": {
						"channelId": "channel-1",
						"title": "Playlist 1"
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
			name: "list playlists json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"playlist-1"}),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list playlists yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"playlist-1"}),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list playlists table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"playlist-1"}),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				p := NewPlaylist(tt.opts...)
				var buf bytes.Buffer
				if err := p.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Playlist.List() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("Playlist.List() output is empty")
				}
			},
		)
	}
}

func TestPlaylist_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert playlist",
			opts: []Option{
				WithTitle("New Playlist"),
				WithDescription("Description"),
				WithPrivacy("public"),
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
		t.Run(
			tt.name, func(t *testing.T) {
				ts := httptest.NewServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							if tt.verify != nil {
								tt.verify(r)
							}
							w.Header().Set("Content-Type", "application/json")
							_, _ = w.Write([]byte(`{"id": "new-playlist-id"}`))
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
				p := NewPlaylist(opts...)
				var buf bytes.Buffer
				if err := p.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Playlist.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestPlaylist_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update playlist",
			opts: []Option{
				WithIds([]string{"playlist-id"}),
				WithTitle("Updated Title"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					parts := r.URL.Query()["part"]
					joined := strings.Join(parts, ",")
					if !strings.Contains(joined, "snippet") || !strings.Contains(
						joined, "status",
					) {
						t.Errorf(
							"expected part to contain snippet and status, got %v", parts,
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update playlist with all fields",
			opts: []Option{
				WithIds([]string{"playlist-id"}),
				WithTitle("Updated Title"),
				WithDescription("Updated Description"),
				WithTags([]string{"tag1", "tag2"}),
				WithLanguage("en"),
				WithPrivacy("private"),
				WithOnBehalfOfContentOwner("owner-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					parts := r.URL.Query()["part"]
					joined := strings.Join(parts, ",")
					if !strings.Contains(joined, "snippet") || !strings.Contains(
						joined, "status",
					) {
						t.Errorf(
							"expected part to contain snippet and status, got %v", parts,
						)
					}
					if r.URL.Query().Get("onBehalfOfContentOwner") != "owner-id" {
						t.Errorf(
							"expected onBehalfOfContentOwner=owner-id, got %s",
							r.URL.Query().Get("onBehalfOfContentOwner"),
						)
					}

					var playlist youtube.Playlist
					if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
						t.Errorf("failed to decode body: %v", err)
					}
					if playlist.Snippet.Title != "Updated Title" {
						t.Errorf(
							"expected title=Updated Title, got %s", playlist.Snippet.Title,
						)
					}
					if playlist.Snippet.Description != "Updated Description" {
						t.Errorf(
							"expected description=Updated Description, got %s",
							playlist.Snippet.Description,
						)
					}
					if playlist.Snippet.DefaultLanguage != "en" {
						t.Errorf(
							"expected language=en, got %s", playlist.Snippet.DefaultLanguage,
						)
					}
					if playlist.Status.PrivacyStatus != "private" {
						t.Errorf(
							"expected privacy=private, got %s", playlist.Status.PrivacyStatus,
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
							{"id": "playlist-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public"}}
						]
					}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "playlist-id", "snippet": {"title": "Updated Title"}}`))
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
				p := NewPlaylist(opts...)
				var buf bytes.Buffer
				if err := p.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Playlist.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestPlaylist_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete playlist",
			opts: []Option{
				WithIds([]string{"playlist-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "delete playlist with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"playlist-id"}),
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
				p := NewPlaylist(opts...)
				var buf bytes.Buffer
				if err := p.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Playlist.Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
