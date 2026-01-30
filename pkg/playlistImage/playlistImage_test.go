// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package playlistImage

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

func TestNewPlaylistImage(t *testing.T) {
	svc := &youtube.Service{}
	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want IPlaylistImage[youtube.PlaylistImage]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"image1", "image2"}),
					WithHeight(1080),
					WithPlaylistId("playlist123"),
					WithType("hero"),
					WithWidth(1920),
					WithFile("/path/to/image.jpg"),
					WithParent("parent123"),
					WithMaxResults(50),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[*].id"),
					WithService(svc),
				},
			},
			want: &PlaylistImage{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[*].id",
				},
				Ids:                           []string{"image1", "image2"},
				Height:                        1080,
				PlaylistId:                    "playlist123",
				Type:                          "hero",
				Width:                         1920,
				File:                          "/path/to/image.jpg",
				Parent:                        "parent123",
				MaxResults:                    50,
				OnBehalfOfContentOwner:        "owner123",
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &PlaylistImage{Fields: &common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &PlaylistImage{
				Fields:     &common.Fields{},
				MaxResults: math.MaxInt64,
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-20),
				},
			},
			want: &PlaylistImage{
				Fields:     &common.Fields{},
				MaxResults: 1,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithPlaylistId(""),
					WithType(""),
					WithFile(""),
					WithParent(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &PlaylistImage{
				Fields:                        &common.Fields{},
				PlaylistId:                    "",
				Type:                          "",
				File:                          "",
				Parent:                        "",
				OnBehalfOfContentOwner:        "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithPlaylistId("myPlaylist"),
					WithType("hero"),
					WithFile("/images/hero.png"),
					WithMaxResults(25),
				},
			},
			want: &PlaylistImage{
				Fields:     &common.Fields{},
				PlaylistId: "myPlaylist",
				Type:       "hero",
				File:       "/images/hero.png",
				MaxResults: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewPlaylistImage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewPlaylistImage() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestPlaylistImage_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get playlist images by parent",
			opts: []Option{
				WithParent("parent-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("parent") != "parent-id" {
					t.Errorf(
						"expected parent=parent-id, got %s", r.URL.Query().Get("parent"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get playlist images with onBehalfOfContentOwner",
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
			name: "get playlist images with onBehalfOfContentOwnerChannel",
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
						{"id": "image-1", "snippet": {"type": "default"}}
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
				pi := NewPlaylistImage(opts...)
				got, err := pi.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("PlaylistImage.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"PlaylistImage.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestPlaylistImage_Get_Pagination(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := 0; i < 20; i++ {
				items[i] = fmt.Sprintf(`{"id": "image-%d"}`, i)
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
				"items": [{"id": "image-20"}, {"id": "image-21"}],
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

	pi := NewPlaylistImage(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := pi.Get()
	if err != nil {
		t.Errorf("PlaylistImage.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("PlaylistImage.Get() got length = %v, want 22", len(got))
	}
}

func TestPlaylistImage_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "image-1",
					"kind": "youtube#playlistImage",
					"snippet": {
						"playlistId": "playlist-1",
						"type": "default"
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
			name: "list playlist images json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list playlist images yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list playlist images table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				pi := NewPlaylistImage(tt.opts...)
				var buf bytes.Buffer
				if err := pi.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistImage.List() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
				if buf.Len() == 0 {
					t.Errorf("PlaylistImage.List() output is empty")
				}
			},
		)
	}
}

func TestPlaylistImage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert playlist image",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithType("default"),
				WithFile("test_image.jpg"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist image with content owner",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithType("default"),
				WithFile("test_image.jpg"),
				WithOnBehalfOfContentOwner("owner-id"),
				WithOnBehalfOfContentOwnerChannel("channel-id"),
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
				if r.URL.Query().Get("onBehalfOfContentOwnerChannel") != "channel-id" {
					t.Errorf(
						"expected onBehalfOfContentOwnerChannel=channel-id, got %s",
						r.URL.Query().Get("onBehalfOfContentOwnerChannel"),
					)
				}
			},
			wantErr: false,
		},
		{
			name: "insert playlist image json output",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithType("default"),
				WithFile("test_image.jpg"),
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
			name: "insert playlist image yaml output",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithType("default"),
				WithFile("test_image.jpg"),
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
			name: "insert playlist image silent output",
			opts: []Option{
				WithPlaylistId("playlist-id"),
				WithType("default"),
				WithFile("test_image.jpg"),
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

	err := os.WriteFile("test_image.jpg", []byte("dummy image content"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer func() {
		_ = os.Remove("test_image.jpg")
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
							_, _ = w.Write([]byte(`{"id": "new-image-id"}`))
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
				pi := NewPlaylistImage(opts...)
				var buf bytes.Buffer
				if err := pi.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistImage.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestPlaylistImage_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update playlist image",
			opts: []Option{
				WithParent("parent-id"), // Need to fetch first to update
				WithMaxResults(1),
				WithType("hero"),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					parts := r.URL.Query()["part"]
					joined := strings.Join(parts, ",")
					if !strings.Contains(joined, "id") || !strings.Contains(
						joined, "kind",
					) || !strings.Contains(joined, "snippet") {
						t.Errorf(
							"expected part to contain id, kind and snippet, got %v", parts,
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update playlist image with content owner",
			opts: []Option{
				WithParent("parent-id"),
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
		{
			name: "update playlist image with file",
			opts: []Option{
				WithParent("parent-id"),
				WithMaxResults(1),
				WithFile("test_image.jpg"),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
				}
			},
			wantErr: false,
		},
	}

	err := os.WriteFile("test_image.jpg", []byte("dummy image content"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer func() {
		_ = os.Remove("test_image.jpg")
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
							if r.Method == "GET" {
								_, _ = w.Write(
									[]byte(`{
						"items": [
							{"id": "image-id", "snippet": {"type": "default"}}
						]
					}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "image-id", "snippet": {"type": "hero"}}`))
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
				pi := NewPlaylistImage(opts...)
				var buf bytes.Buffer
				if err := pi.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistImage.Update() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestPlaylistImage_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete playlist image",
			opts: []Option{
				WithIds([]string{"image-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				if r.URL.Query().Get("id") != "image-id" {
					t.Errorf("expected id=image-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantErr: false,
		},
		{
			name: "delete playlist image with content owner",
			opts: []Option{
				WithIds([]string{"image-id"}),
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
				pi := NewPlaylistImage(opts...)
				var buf bytes.Buffer
				if err := pi.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"PlaylistImage.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
