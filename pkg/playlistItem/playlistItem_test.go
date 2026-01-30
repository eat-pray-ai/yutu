// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistItem

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

func TestNewPlaylistItem(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IPlaylistItem[youtube.PlaylistItem]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"item1", "item2"}),
					WithTitle("Test Item"),
					WithDescription("Test item description"),
					WithKind("video"),
					WithKVideoId("video123"),
					WithKChannelId("channel123"),
					WithKPlaylistId("playlist123"),
					WithVideoId("video456"),
					WithPlaylistId("playlist456"),
					WithChannelId("channel456"),
					WithPrivacy("public"),
					WithMaxResults(50),
					WithOnBehalfOfContentOwner("owner123"),
					WithParts([]string{"snippet", "status"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &PlaylistItem{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"snippet", "status"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				Ids:                    []string{"item1", "item2"},
				Title:                  "Test Item",
				Description:            "Test item description",
				Kind:                   "video",
				KVideoId:               "video123",
				KChannelId:             "channel123",
				KPlaylistId:            "playlist123",
				VideoId:                "video456",
				PlaylistId:             "playlist456",
				ChannelId:              "channel456",
				Privacy:                "public",
				MaxResults:             50,
				OnBehalfOfContentOwner: "owner123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &PlaylistItem{Fields: &common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &PlaylistItem{
				Fields:     &common.Fields{},
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-15),
				},
			},
			want: &PlaylistItem{
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
					WithKind(""),
					WithKVideoId(""),
					WithKChannelId(""),
					WithKPlaylistId(""),
					WithVideoId(""),
					WithPlaylistId(""),
					WithChannelId(""),
					WithPrivacy(""),
					WithOnBehalfOfContentOwner(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &PlaylistItem{
				Fields:                 &common.Fields{Output: "", Jsonpath: ""},
				Title:                  "",
				Description:            "",
				Kind:                   "",
				KVideoId:               "",
				KChannelId:             "",
				KPlaylistId:            "",
				VideoId:                "",
				PlaylistId:             "",
				ChannelId:              "",
				Privacy:                "",
				OnBehalfOfContentOwner: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Video"),
					WithKind("video"),
					WithKVideoId("myVideo123"),
					WithPlaylistId("myPlaylist"),
					WithMaxResults(25),
					WithParts([]string{"id"}),
				},
			},
			want: &PlaylistItem{
				Fields:     &common.Fields{Parts: []string{"id"}},
				Title:      "My Video",
				Kind:       "video",
				KVideoId:   "myVideo123",
				PlaylistId: "myPlaylist",
				MaxResults: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylistItem(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewPlaylistItem() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestPlaylistItem_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get playlist items by id",
			opts: []Option{
				WithIds([]string{"item-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "item-id" {
					t.Errorf("expected id=item-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get playlist items by playlistId",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("playlistId") != "playlist-id" {
					t.Errorf(
						"expected playlistId=playlist-id, got %s",
						r.URL.Query().Get("playlistId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get playlist items by videoId",
			opts: []Option{
				WithVideoId("video-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("videoId") != "video-id" {
					t.Errorf(
						"expected videoId=video-id, got %s", r.URL.Query().Get("videoId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get playlist items with onBehalfOfContentOwner",
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
						{"id": "item-1", "snippet": {"title": "Item 1", "resourceId": {"kind": "youtube#video", "videoId": "video-1"}}}
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
				pi := NewPlaylistItem(opts...)
				got, err := pi.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("PlaylistItem.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"PlaylistItem.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestPlaylistItem_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := 0; i < 20; i++ {
				items[i] = fmt.Sprintf(`{"id": "item-%d"}`, i)
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
				"items": [{"id": "item-20"}, {"id": "item-21"}],
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

	pi := NewPlaylistItem(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := pi.Get()
	if err != nil {
		t.Errorf("PlaylistItem.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("PlaylistItem.Get() got length = %v, want 22", len(got))
	}
}

func TestPlaylistItem_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "item-1",
					"snippet": {
						"title": "Item 1",
						"resourceId": {
							"kind": "youtube#video",
							"videoId": "video-1"
						}
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
			name: "list playlist items json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"item-1"}),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list playlist items yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"item-1"}),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list playlist items table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"item-1"}),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				pi := NewPlaylistItem(tt.opts...)
				var buf bytes.Buffer
				if err := pi.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf("PlaylistItem.List() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("PlaylistItem.List() output is empty")
				}
			},
		)
	}
}

func TestPlaylistItem_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert playlist item video",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("video"),
				WithKVideoId("video-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist item channel",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("channel"),
				WithKChannelId("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist item playlist",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("playlist"),
				WithKPlaylistId("playlist-id-2"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist item with content owner and privacy",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("video"),
				WithKVideoId("video-id"),
				WithOnBehalfOfContentOwner("owner-id"),
				WithPrivacy("private"),
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
		{
			name: "insert playlist item json output",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("video"),
				WithKVideoId("video-id"),
				WithOutput("json"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist item yaml output",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("video"),
				WithKVideoId("video-id"),
				WithOutput("yaml"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist item silent output",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithKind("video"),
				WithKVideoId("video-id"),
				WithOutput("silent"),
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
							_, _ = w.Write([]byte(`{"id": "new-item-id"}`))
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
				pi := NewPlaylistItem(opts...)
				var buf bytes.Buffer
				if err := pi.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistItem.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestPlaylistItem_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update playlist item",
			opts: []Option{
				WithIds([]string{"item-id"}),
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
			name: "update playlist item description and privacy",
			opts: []Option{
				WithIds([]string{"item-id"}),
				WithDescription("Updated Description"),
				WithPrivacy("private"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("id") != "" {
					}
				}
			},
			wantErr: false,
		},

		{
			name: "update playlist item with content owner",
			opts: []Option{
				WithIds([]string{"item-id"}),
				WithTitle("Updated Title"),
				WithOnBehalfOfContentOwner("owner-id"),
				WithMaxResults(1),
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
							{"id": "item-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public"}}
						]
					}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "item-id", "snippet": {"title": "Updated Title"}}`))
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
				pi := NewPlaylistItem(opts...)
				var buf bytes.Buffer
				if err := pi.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistItem.Update() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestPlaylistItem_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete playlist item",
			opts: []Option{
				WithIds([]string{"item-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "delete playlist item with content owner",
			opts: []Option{
				WithIds([]string{"item-id"}),
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
				pi := NewPlaylistItem(opts...)
				var buf bytes.Buffer
				if err := pi.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistItem.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
