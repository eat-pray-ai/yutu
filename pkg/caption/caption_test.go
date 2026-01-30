// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package caption

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestCaption_GetFields(t *testing.T) {
	c := NewCaption().(*Caption)
	if c.GetFields() == nil {
		t.Error("Caption.GetFields() is nil")
	}
}

func TestNewCaption(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
	isAutoSyncedTrue := true
	isAutoSyncedFalse := false
	isCCTrue := true
	isCCFalse := false
	isDraftTrue := true
	isDraftFalse := false
	isEasyReaderTrue := true
	isEasyReaderFalse := false
	isLargeTrue := true
	isLargeFalse := false

	tests := []struct {
		name string
		args args
		want ICaption[youtube.Caption]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"caption1", "caption2"}),
					WithFile("/path/to/file.srt"),
					WithAudioTrackType("primary"),
					WithIsAutoSynced(&isAutoSyncedTrue),
					WithIsCC(&isCCTrue),
					WithIsDraft(&isDraftTrue),
					WithIsEasyReader(&isEasyReaderTrue),
					WithIsLarge(&isLargeTrue),
					WithLanguage("en"),
					WithName("English Captions"),
					WithTrackKind("standard"),
					WithOnBehalfOf("channel123"),
					WithOnBehalfOfContentOwner("owner123"),
					WithVideoId("video123"),
					WithTfmt("srt"),
					WithTlang("es"),
					WithService(svc),
				},
			},
			want: &Caption{
				Fields:                 &common.Fields{Service: svc},
				Ids:                    []string{"caption1", "caption2"},
				File:                   "/path/to/file.srt",
				AudioTrackType:         "primary",
				IsAutoSynced:           &isAutoSyncedTrue,
				IsCC:                   &isCCTrue,
				IsDraft:                &isDraftTrue,
				IsEasyReader:           &isEasyReaderTrue,
				IsLarge:                &isLargeTrue,
				Language:               "en",
				Name:                   "English Captions",
				TrackKind:              "standard",
				OnBehalfOf:             "channel123",
				OnBehalfOfContentOwner: "owner123",
				VideoId:                "video123",
				Tfmt:                   "srt",
				Tlang:                  "es",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Caption{Fields: &common.Fields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithIsAutoSynced(nil),
					WithIsCC(nil),
					WithIsDraft(nil),
					WithIsEasyReader(nil),
					WithIsLarge(nil),
				},
			},
			want: &Caption{Fields: &common.Fields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithIsAutoSynced(&isAutoSyncedFalse),
					WithIsCC(&isCCFalse),
					WithIsDraft(&isDraftFalse),
					WithIsEasyReader(&isEasyReaderFalse),
					WithIsLarge(&isLargeFalse),
				},
			},
			want: &Caption{
				Fields:       &common.Fields{},
				IsAutoSynced: &isAutoSyncedFalse,
				IsCC:         &isCCFalse,
				IsDraft:      &isDraftFalse,
				IsEasyReader: &isEasyReaderFalse,
				IsLarge:      &isLargeFalse,
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithFile(""),
					WithAudioTrackType(""),
					WithLanguage(""),
					WithName(""),
					WithTrackKind(""),
					WithOnBehalfOf(""),
					WithOnBehalfOfContentOwner(""),
					WithVideoId(""),
					WithTfmt(""),
					WithTlang(""),
				},
			},
			want: &Caption{
				Fields:                 &common.Fields{},
				File:                   "",
				AudioTrackType:         "",
				Language:               "",
				Name:                   "",
				TrackKind:              "",
				OnBehalfOf:             "",
				OnBehalfOfContentOwner: "",
				VideoId:                "",
				Tfmt:                   "",
				Tlang:                  "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithIds([]string{"caption1"}),
					WithLanguage("fr"),
					WithVideoId("video456"),
					WithIsCC(&isCCTrue),
				},
			},
			want: &Caption{
				Fields:   &common.Fields{},
				Ids:      []string{"caption1"},
				Language: "fr",
				VideoId:  "video456",
				IsCC:     &isCCTrue,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewCaption(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewCaption() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestCaption_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get captions with videoId",
			opts: []Option{
				WithVideoId("video-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("videoId") != "video-id" {
					t.Errorf(
						"expected videoId=video-id, got %s", r.URL.Query().Get("videoId"),
					)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get captions with ids",
			opts: []Option{
				WithVideoId("video-id"),
				WithIds([]string{"id1", "id2"}),
			},
			verify: func(r *http.Request) {
				ids := r.URL.Query()["id"]
				if len(ids) != 2 || ids[0] != "id1" || ids[1] != "id2" {
					t.Errorf("expected id=[id1 id2], got %v", ids)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get captions with onBehalfOf",
			opts: []Option{
				WithVideoId("video-id"),
				WithOnBehalfOf("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOf") != "channel-id" {
					t.Errorf(
						"expected onBehalfOf=channel-id, got %s",
						r.URL.Query().Get("onBehalfOf"),
					)
				}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "get captions with onBehalfOfContentOwner",
			opts: []Option{
				WithVideoId("video-id"),
				WithOnBehalfOfContentOwner("content-owner"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "content-owner" {
					t.Errorf(
						"expected onBehalfOfContentOwner=content-owner, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
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
						{"id": "1", "snippet": {"videoId": "video-id"}},
						{"id": "2", "snippet": {"videoId": "video-id"}}
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
				c := NewCaption(opts...)
				got, err := c.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Caption.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"Caption.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestCaption_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "caption-1",
					"snippet": {
						"videoId": "video-1",
						"name": "English",
						"language": "en"
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
			name: "list captions json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithVideoId("video-1"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list captions yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithVideoId("video-1"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list captions table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithVideoId("video-1"),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := NewCaption(tt.opts...)
				var buf bytes.Buffer
				if err := c.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Caption.List() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("Caption.List() output is empty")
				}
			},
		)
	}
}

func TestCaption_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert caption",
			opts: []Option{
				WithVideoId("video-id"),
				WithLanguage("en"),
				WithName("English"),
				WithFile("test_caption.srt"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "insert caption with onBehalfOf",
			opts: []Option{
				WithVideoId("video-id"),
				WithFile("test_caption.srt"),
				WithOnBehalfOf("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOf") != "channel-id" {
					t.Errorf(
						"expected onBehalfOf=channel-id, got %s",
						r.URL.Query().Get("onBehalfOf"),
					)
				}
			},
			wantErr: false,
		},
		{
			name: "insert caption with onBehalfOfContentOwner",
			opts: []Option{
				WithVideoId("video-id"),
				WithFile("test_caption.srt"),
				WithOnBehalfOfContentOwner("content-owner"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "content-owner" {
					t.Errorf(
						"expected onBehalfOfContentOwner=content-owner, got %s",
						r.URL.Query().Get("onBehalfOfContentOwner"),
					)
				}
			},
			wantErr: false,
		},
	}

	err := os.WriteFile(
		"test_caption.srt", []byte("1\n00:00:01,000 --> 00:00:04,000\nHello"), 0644,
	)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
	}
	defer func() {
		_ = os.Remove("test_caption.srt")
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
							_, _ = w.Write([]byte(`{"id": "new-caption-id", "snippet": {"videoId": "video-id"}}`))
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
				c := NewCaption(opts...)
				var buf bytes.Buffer
				if err := c.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Caption.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestCaption_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update caption",
			opts: []Option{
				WithVideoId("video-id"),
				WithIds([]string{"caption-id"}),
				WithLanguage("es"),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "snippet" {
						t.Errorf("expected part=snippet, got %s", r.URL.Query().Get("part"))
					}
				} else if r.Method == "GET" {
					if r.URL.Query().Get("videoId") != "video-id" {
						t.Errorf(
							"expected videoId=video-id, got %s", r.URL.Query().Get("videoId"),
						)
					}
				} else {
					t.Errorf("unexpected method %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "update caption with onBehalfOf",
			opts: []Option{
				WithVideoId("video-id"),
				WithIds([]string{"caption-id"}),
				WithOnBehalfOf("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("onBehalfOf") != "channel-id" {
						t.Errorf(
							"expected onBehalfOf=channel-id, got %s",
							r.URL.Query().Get("onBehalfOf"),
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update caption with onBehalfOfContentOwner",
			opts: []Option{
				WithVideoId("video-id"),
				WithIds([]string{"caption-id"}),
				WithOnBehalfOfContentOwner("content-owner"),
			},
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("onBehalfOfContentOwner") != "content-owner" {
						t.Errorf(
							"expected onBehalfOfContentOwner=content-owner, got %s",
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
							{"id": "caption-id", "snippet": {"videoId": "video-id", "language": "en"}}
						]
					}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "caption-id", "snippet": {"videoId": "video-id", "language": "es"}}`))
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
				c := NewCaption(opts...)
				var buf bytes.Buffer
				if err := c.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Caption.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestCaption_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete caption",
			opts: []Option{
				WithIds([]string{"caption-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "delete caption with onBehalfOf",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithOnBehalfOf("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOf") != "channel-id" {
					t.Errorf(
						"expected onBehalfOf=channel-id, got %s",
						r.URL.Query().Get("onBehalfOf"),
					)
				}
			},
			wantErr: false,
		},
		{
			name: "delete caption with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithOnBehalfOfContentOwner("content-owner"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "content-owner" {
					t.Errorf(
						"expected onBehalfOfContentOwner=content-owner, got %s",
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
				c := NewCaption(opts...)
				var buf bytes.Buffer
				if err := c.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Caption.Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestCaption_Download(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "download caption",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithFile("downloaded.srt"),
			},
			verify: func(r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "download caption with tfmt",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithFile("downloaded.srt"),
				WithTfmt("vtt"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("tfmt") != "vtt" {
					t.Errorf("expected tfmt=vtt, got %s", r.URL.Query().Get("tfmt"))
				}
			},
			wantErr: false,
		},
		{
			name: "download caption with tlang",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithFile("downloaded.srt"),
				WithTlang("es"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("tlang") != "es" {
					t.Errorf("expected tlang=es, got %s", r.URL.Query().Get("tlang"))
				}
			},
			wantErr: false,
		},
		{
			name: "download caption with onBehalfOf",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithFile("downloaded.srt"),
				WithOnBehalfOf("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOf") != "channel-id" {
					t.Errorf(
						"expected onBehalfOf=channel-id, got %s",
						r.URL.Query().Get("onBehalfOf"),
					)
				}
			},
			wantErr: false,
		},
		{
			name: "download caption with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"caption-id"}),
				WithFile("downloaded.srt"),
				WithOnBehalfOfContentOwner("content-owner"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("onBehalfOfContentOwner") != "content-owner" {
					t.Errorf(
						"expected onBehalfOfContentOwner=content-owner, got %s",
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
							_, _ = w.Write([]byte("1\n00:00:01,000 --> 00:00:04,000\nDownloaded Caption"))
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
				c := NewCaption(opts...)
				var buf bytes.Buffer
				if err := c.Download(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Caption.Download() error = %v, wantErr %v", err, tt.wantErr)
				}
				defer func() {
					_ = os.Remove("downloaded.srt")
				}()
			},
		)
	}
}

func TestCaption_Insert_Error(t *testing.T) {
	// Setup pkg.Root for file tests
	tmpDir := t.TempDir()
	f, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	oldRoot := pkg.Root
	pkg.Root = f
	defer func() { pkg.Root = oldRoot }()
	defer func(f *os.Root) {
		_ = f.Close()
	}(f)

	svc, _ := youtube.NewService(context.Background(), option.WithAPIKey("test"))

	// Test: File open error
	c := NewCaption(WithFile("non_existent.srt"), WithService(svc))
	if err := c.Insert(&bytes.Buffer{}); err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// Test: API error
	// Create dummy file
	if err := os.WriteFile(
		tmpDir+"/test.srt", []byte("content"), 0644,
	); err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		),
	)
	defer ts.Close()

	svc, _ = youtube.NewService(
		context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test"),
	)
	c = NewCaption(WithFile("test.srt"), WithService(svc))
	if err := c.Insert(&bytes.Buffer{}); err == nil {
		t.Error("expected error from API, got nil")
	}
}

func TestCaption_Update_Error(t *testing.T) {
	// Setup pkg.Root
	tmpDir := t.TempDir()
	f, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	oldRoot := pkg.Root
	pkg.Root = f
	defer func() { pkg.Root = oldRoot }()
	defer func(f *os.Root) {
		_ = f.Close()
	}(f)

	// Test: File open error (when file is specified)
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"items": [{"id": "1", "snippet": {"videoId": "v1"}}]}`))
			},
		),
	)
	defer ts.Close()
	svc, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test"),
	)

	c := NewCaption(
		WithFile("non_existent.srt"), WithService(svc), WithVideoId("v1"),
	)
	if err := c.Update(&bytes.Buffer{}); err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// Test: Get error
	tsError := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		),
	)
	defer tsError.Close()
	svcError, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(tsError.URL),
		option.WithAPIKey("test"),
	)
	c = NewCaption(WithService(svcError), WithVideoId("v1"))
	if err := c.Update(&bytes.Buffer{}); err == nil {
		t.Error("expected error from Get, got nil")
	}

	// Test: No caption found
	tsEmpty := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"items": []}`))
			},
		),
	)
	defer tsEmpty.Close()
	svcEmpty, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(tsEmpty.URL),
		option.WithAPIKey("test"),
	)
	c = NewCaption(WithService(svcEmpty), WithVideoId("v1"))
	if err := c.Update(&bytes.Buffer{}); !errors.Is(err, errGetCaption) {
		t.Errorf("expected errGetCaption, got %v", err)
	}

	// Test: API Update Error
	tsUpdateErr := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"items": [{"id": "1", "snippet": {"videoId": "v1"}}]}`))
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
			},
		),
	)
	defer tsUpdateErr.Close()
	svcUpdateErr, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(tsUpdateErr.URL),
		option.WithAPIKey("test"),
	)
	c = NewCaption(WithService(svcUpdateErr), WithVideoId("v1"))
	if err := c.Update(&bytes.Buffer{}); err == nil {
		t.Error("expected error from API Update, got nil")
	}
}

func TestCaption_Delete_Error(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		),
	)
	defer ts.Close()
	svc, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test"),
	)

	c := NewCaption(WithService(svc), WithIds([]string{"id1"}))
	if err := c.Delete(&bytes.Buffer{}); err == nil {
		t.Error("expected error from Delete, got nil")
	}
}

func TestCaption_Download_Error(t *testing.T) {
	// Test API Error
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
		),
	)
	defer ts.Close()
	svc, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(ts.URL), option.WithAPIKey("test"),
	)

	c := NewCaption(
		WithService(svc), WithIds([]string{"id1"}), WithFile("out.srt"),
	)
	if err := c.Download(&bytes.Buffer{}); err == nil {
		t.Error("expected error from Download API, got nil")
	}

	// Test File Creation Error
	ts2 := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("content"))
			},
		),
	)
	defer ts2.Close()
	svc2, _ := youtube.NewService(
		context.Background(), option.WithEndpoint(ts2.URL),
		option.WithAPIKey("test"),
	)
	// Use a directory as file path to trigger error.
	tmpDir := t.TempDir()
	c = NewCaption(WithService(svc2), WithIds([]string{"id1"}), WithFile(tmpDir))
	if err := c.Download(&bytes.Buffer{}); err == nil {
		t.Error("expected error for file creation, got nil")
	}
}
