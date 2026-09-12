// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"encoding/json/v2"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestVideo_UpdateClearableMetadata(t *testing.T) {
	for _, tt := range []struct {
		name              string
		input             string
		privacy           string
		wantDescription   string
		wantLanguage      string
		wantPublishAt     string
		wantRecordingDate string
		recordingPart     bool
	}{
		{name: "omitted preserves", input: `{}`, privacy: "private", wantDescription: "Original", wantLanguage: "en", wantPublishAt: "2027-01-01T00:00:00Z"},
		{name: "empty clears", input: `{"description":"","language":"","publish_at":"","recording_date":""}`, privacy: "private", recordingPart: true},
		{name: "clear does not require private", input: `{"publish_at":""}`, privacy: "public", wantDescription: "Original", wantLanguage: "en"},
		{name: "nonempty replaces", input: `{"description":"New","language":"ja","publish_at":"2028-01-01T00:00:00Z","recording_date":"2020-01-01T00:00:00Z"}`, privacy: "private", wantDescription: "New", wantLanguage: "ja", wantPublishAt: "2028-01-01T00:00:00Z", wantRecordingDate: "2020-01-01T00:00:00Z", recordingPart: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var body map[string]map[string]any
			var parts string
			svc := common.NewTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodGet {
					_, _ = io.WriteString(w, `{"items":[{"id":"video-id","snippet":{"title":"Title","categoryId":"22","description":"Original","defaultLanguage":"en"},"status":{"privacyStatus":"`+tt.privacy+`","publishAt":"2027-01-01T00:00:00Z"}}]}`)
					return
				}
				parts = r.URL.Query().Get("part")
				var raw struct {
					Snippet          map[string]any `json:"snippet"`
					Status           map[string]any `json:"status"`
					RecordingDetails map[string]any `json:"recordingDetails"`
				}
				if err := json.UnmarshalRead(r.Body, &raw); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				body = map[string]map[string]any{"snippet": raw.Snippet, "status": raw.Status, "recordingDetails": raw.RecordingDetails}
				_, _ = io.WriteString(w, `{"id":"video-id"}`)
			}))
			var input Video
			if err := json.Unmarshal([]byte(tt.input), &input); err != nil {
				t.Fatal(err)
			}
			input.Service, input.Ids, input.MaxResults, input.Output = svc, []string{"video-id"}, 1, "silent"
			if err := input.Update(io.Discard); err != nil {
				t.Fatal(err)
			}
			if body == nil {
				t.Fatal("missing update request")
			}
			for _, field := range []struct{ part, key, want string }{
				{"snippet", "description", tt.wantDescription},
				{"snippet", "defaultLanguage", tt.wantLanguage},
				{"status", "publishAt", tt.wantPublishAt},
				{"recordingDetails", "recordingDate", tt.wantRecordingDate},
			} {
				got, present := body[field.part][field.key]
				if field.want == "" {
					if present {
						t.Errorf("%s.%s should be omitted for clearing, got %v", field.part, field.key, got)
					}
				} else if got != field.want {
					t.Errorf("%s.%s = %v, want %q", field.part, field.key, got, field.want)
				}
			}
			if strings.Contains(parts, "recordingDetails") != tt.recordingPart {
				t.Errorf("unexpected update parts: %s", parts)
			}
			if tt.recordingPart && body["recordingDetails"] == nil {
				t.Error("recordingDetails object missing")
			}
		})
	}
}

func TestVideo_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get videos by id",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "video-id" {
					t.Errorf("expected id=video-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get videos myRating",
			opts: []Option{
				WithRating("like"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("myRating") != "like" {
					t.Errorf(
						"expected myRating=like, got %s", r.URL.Query().Get("myRating"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get videos chart",
			opts: []Option{
				WithChart("mostPopular"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("chart") != "mostPopular" {
					t.Errorf(
						"expected chart=mostPopular, got %s", r.URL.Query().Get("chart"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get videos with regionCode, maxHeight, maxWidth, onBehalfOfContentOwner",
			opts: []Option{
				WithRegionCode("US"),
				WithMaxHeight(1080),
				WithMaxWidth(1920),
				WithOnBehalfOfContentOwner("owner-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("regionCode") != "US" {
					t.Errorf(
						"expected regionCode=US, got %s", r.URL.Query().Get("regionCode"),
					)
				}
				if r.URL.Query().Get("maxHeight") != "1080" {
					t.Errorf(
						"expected maxHeight=1080, got %s", r.URL.Query().Get("maxHeight"),
					)
				}
				if r.URL.Query().Get("maxWidth") != "1920" {
					t.Errorf(
						"expected maxWidth=1920, got %s", r.URL.Query().Get("maxWidth"),
					)
				}
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
			name: "get videos with hl, locale, categoryId",
			opts: []Option{
				WithHl("en"),
				WithLocale("en_US"),
				WithCategory("10"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "en" {
					t.Errorf("expected hl=en, got %s", r.URL.Query().Get("hl"))
				}
				if r.URL.Query().Get("locale") != "en_US" {
					t.Errorf("expected locale=en_US, got %s", r.URL.Query().Get("locale"))
				}
				if r.URL.Query().Get("videoCategoryId") != "10" {
					t.Errorf(
						"expected videoCategoryId=10, got %s",
						r.URL.Query().Get("videoCategoryId"),
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
						{"id": "video-1", "snippet": {"title": "Video 1"}}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				got, err := v.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Video.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("Video.Get() got length = %v, want %v", len(got), tt.wantLen)
				}
			},
		)
	}
}

func TestVideo_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("video"))

	v := NewVideo(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := v.Get()
	if err != nil {
		t.Errorf("Video.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Video.Get() got length = %v, want 22", len(got))
	}
}

func TestVideo_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "video-1",
				"snippet": {
					"channelId": "channel-1",
					"title": "Video 1"
				},
				"statistics": {
					"viewCount": "100"
				}
			}
		]
	}`

	common.RunListTest(
		t, mockResponse,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			v := NewVideo(
				WithService(svc),
				WithOutput(output),
				WithIds([]string{"video-1"}),
				WithMaxResults(5),
			)
			return v.List
		},
	)
}

func TestVideo_List_NilFields(t *testing.T) {
	mockResponse := `{
		"items": [
			{"id": "video-1"},
			{"id": "video-2", "snippet": {"title": "Video 2", "channelId": "ch-2"}},
			{"id": "video-3", "statistics": {"viewCount": "500"}}
		]
	}`

	svc := common.NewTestService(
		t, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(mockResponse))
			},
		),
	)

	v := NewVideo(
		WithService(svc),
		WithOutput("table"),
		WithIds([]string{"video-1"}),
		WithMaxResults(5),
	)
	var buf bytes.Buffer
	if err := v.List(&buf); err != nil {
		t.Errorf("List() error = %v", err)
	}
	if buf.Len() == 0 {
		t.Error("List() output is empty")
	}
}

func TestVideo_Insert(t *testing.T) {
	autoLevelsTrue := true
	notifySubscribersTrue := true
	stabilizeTrue := true
	forKidsTrue := true
	embeddableTrue := true
	containsSyntheticMediaTrue := true
	publicStatsViewableTrue := true

	decodeMultipartVideo := func(t *testing.T, r *http.Request) *youtube.Video {
		t.Helper()
		ct := r.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(ct)
		if err != nil {
			t.Fatalf("failed to parse Content-Type %q: %v", ct, err)
		}
		if !strings.HasPrefix(mediaType, "multipart/") {
			t.Fatalf("expected multipart content type, got %s", mediaType)
		}
		boundary, ok := params["boundary"]
		if !ok || boundary == "" {
			t.Fatalf("missing multipart boundary in Content-Type %q", ct)
		}

		mr := multipart.NewReader(r.Body, boundary)
		part, err := mr.NextPart()
		if err != nil {
			t.Fatalf("failed to read first multipart part: %v", err)
		}
		defer func() { _ = part.Close() }()

		var body youtube.Video
		if err := json.UnmarshalRead(part, &body); err != nil {
			t.Fatalf("failed to decode video from request body: %v", err)
		}
		return &body
	}

	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert explicit false status fields",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithEmbeddable(new(false)),
				WithPublicStatsViewable(new(false)),
			},
			verify: func(r *http.Request) {
				_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if err != nil {
					t.Error(err)
					return
				}
				part, err := multipart.NewReader(r.Body, params["boundary"]).NextPart()
				if err != nil {
					t.Error(err)
					return
				}
				defer func() { _ = part.Close() }()
				var body struct {
					Status map[string]any `json:"status"`
				}
				if err := json.UnmarshalRead(part, &body); err != nil {
					t.Error(err)
					return
				}
				for _, field := range []string{"embeddable", "publicStatsViewable"} {
					if value, present := body.Status[field]; !present || value != false {
						t.Errorf("status.%s must be present and false, got %v (present=%v)", field, value, present)
					}
				}
			},
		},
		{
			name: "insert video",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithTitle("New Video"),
				WithPrivacy("public"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				body := decodeMultipartVideo(t, r)
				if body.Snippet == nil || body.Snippet.Title != "New Video" {
					t.Errorf("expected snippet.title=New Video, got %+v", body.Snippet)
				}
				if body.Status == nil || body.Status.PrivacyStatus != "public" {
					t.Errorf("expected status.privacyStatus=public, got %+v", body.Status)
				}
			},
			wantErr: false,
		},
		{
			name: "insert video with options",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithTitle("New Video"),
				WithPrivacy("public"),
				WithAutoLevels(&autoLevelsTrue),
				WithNotifySubscribers(&notifySubscribersTrue),
				WithStabilize(&stabilizeTrue),
				WithOnBehalfOfContentOwner("owner-id"),
				WithOnBehalfOfContentOwnerChannel("channel-id"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("autoLevels") != "true" {
					t.Errorf(
						"expected autoLevels=true, got %s", r.URL.Query().Get("autoLevels"),
					)
				}
				if r.URL.Query().Get("notifySubscribers") != "true" {
					t.Errorf(
						"expected notifySubscribers=true, got %s",
						r.URL.Query().Get("notifySubscribers"),
					)
				}
				if r.URL.Query().Get("stabilize") != "true" {
					t.Errorf(
						"expected stabilize=true, got %s", r.URL.Query().Get("stabilize"),
					)
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
				body := decodeMultipartVideo(t, r)
				if body.Snippet == nil || body.Snippet.Title != "New Video" {
					t.Errorf("expected snippet.title=New Video, got %+v", body.Snippet)
				}
				if body.Status == nil || body.Status.PrivacyStatus != "public" {
					t.Errorf("expected status.privacyStatus=public, got %+v", body.Status)
				}
			},
			wantErr: false,
		},
		{
			name: "insert video with recording date and boolean options",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithTitle("Dated Video"),
				WithPrivacy("private"),
				WithRecordingDate(new("2024-06-15T10:00:00Z")),
				WithForKids(&forKidsTrue),
				WithEmbeddable(&embeddableTrue),
				WithContainsSyntheticMedia(&containsSyntheticMediaTrue),
				WithPublicStatsViewable(&publicStatsViewableTrue),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if !strings.Contains(r.URL.Query().Get("part"), "recordingDetails") {
					t.Errorf(
						"expected part to contain recordingDetails, got %s",
						r.URL.Query().Get("part"),
					)
				}
				body := decodeMultipartVideo(t, r)
				if body.RecordingDetails == nil || body.RecordingDetails.RecordingDate != "2024-06-15T10:00:00Z" {
					t.Errorf(
						"expected recordingDetails.recordingDate=2024-06-15T10:00:00Z, got %+v",
						body.RecordingDetails,
					)
				}
				if body.Status == nil || !body.Status.SelfDeclaredMadeForKids {
					t.Errorf("expected status.selfDeclaredMadeForKids=true")
				}
				if body.Status == nil || !body.Status.Embeddable {
					t.Errorf("expected status.embeddable=true")
				}
				if body.Status == nil || !body.Status.ContainsSyntheticMedia {
					t.Errorf("expected status.containsSyntheticMedia=true")
				}
				if body.Status == nil || !body.Status.PublicStatsViewable {
					t.Errorf("expected status.publicStatsViewable=true")
				}
			},
			wantErr: false,
		},
		{
			name: "insert video with empty title uses filename",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithPrivacy("public"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					return
				}
				body := decodeMultipartVideo(t, r)
				if body.Snippet == nil || body.Snippet.Title != "test_video" {
					t.Errorf(
						"expected snippet.title=test_video (from filename), got %+v",
						body.Snippet,
					)
				}
			},
			wantErr: false,
		},
		{
			name: "insert video with existing yutu tag",
			opts: []Option{
				WithFile("test_video.mp4"),
				WithTitle("Tagged Video"),
				WithTags([]string{"yutu🐰", "existing"}),
				WithPrivacy("public"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					return
				}
				body := decodeMultipartVideo(t, r)
				count := 0
				for _, tag := range body.Snippet.Tags {
					if tag == "yutu🐰" {
						count++
					}
				}
				if count != 1 {
					t.Errorf(
						"expected exactly 1 yutu🐰 tag, got %d in %v", count,
						body.Snippet.Tags,
					)
				}
			},
			wantErr: false,
		},
	}

	tmpDir := t.TempDir()
	root, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatalf("failed to open root: %v", err)
	}
	oldRoot := pkg.Root
	pkg.Root = root
	defer func() { pkg.Root = oldRoot }()
	defer func() { _ = root.Close() }()

	err = os.WriteFile(
		tmpDir+"/test_video.mp4", []byte("dummy video content"), 0644,
	)
	if err != nil {
		t.Fatalf("failed to create dummy file: %v", err)
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
							_, _ = w.Write([]byte(`{"id": "new-video-id", "snippet": {"title": "New Video"}, "status": {"privacyStatus": "public"}}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_Insert_FileError(t *testing.T) {
	tmpDir := t.TempDir()
	root, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatalf("failed to open root: %v", err)
	}
	oldRoot := pkg.Root
	pkg.Root = root
	defer func() { pkg.Root = oldRoot }()
	defer func() { _ = root.Close() }()

	svc := common.NewTestService(
		t, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	v := NewVideo(
		WithService(svc),
		WithFile("non_existent.mp4"),
	)
	var buf bytes.Buffer
	if err := v.Insert(&buf); err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestVideo_Update(t *testing.T) {
	embeddableTrue := true
	containsSyntheticMediaTrue := true

	tests := []struct {
		name        string
		opts        []Option
		getResponse string
		verify      func(*http.Request)
		wantErr     bool
	}{
		{
			name: "update video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithTitle("Updated Title"),
				WithDescription(new("Updated Description")),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if r.URL.Query().Get("part") != "snippet,status" {
						t.Errorf(
							"expected part=snippet,status, got %s", r.URL.Query().Get("part"),
						)
					}

					var body youtube.Video
					if err := json.UnmarshalRead(r.Body, &body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if body.Id != "video-id" {
						t.Errorf("expected id=video-id, got %s", body.Id)
					}
					if body.Snippet == nil || body.Snippet.Title != "Updated Title" {
						t.Errorf(
							"expected snippet.title=Updated Title, got %+v", body.Snippet,
						)
					}
					if body.Snippet.Description != "Updated Description" {
						t.Errorf(
							"expected snippet.description=Updated Description, got %s",
							body.Snippet.Description,
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update video with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithTitle("Updated Title"),
				WithMaxResults(1),
				WithOnBehalfOfContentOwner("owner-id"),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}}]}`,
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
			name: "update video with tags, language, license, categoryId, privacy",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithTags([]string{"new-tag"}),
				WithLanguage(new("ja")),
				WithLicense("creativeCommon"),
				WithCategory("22"),
				WithPrivacy("unlisted"),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public", "license": "youtube", "embeddable": true, "publicStatsViewable": true}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					var body youtube.Video
					if err := json.UnmarshalRead(r.Body, &body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if body.Snippet.DefaultLanguage != "ja" {
						t.Errorf(
							"expected defaultLanguage=ja, got %s", body.Snippet.DefaultLanguage,
						)
					}
					if body.Status.License != "creativeCommon" {
						t.Errorf(
							"expected license=creativeCommon, got %s", body.Status.License,
						)
					}
					if body.Snippet.CategoryId != "22" {
						t.Errorf("expected categoryId=22, got %s", body.Snippet.CategoryId)
					}
					if body.Status.PrivacyStatus != "unlisted" {
						t.Errorf(
							"expected privacyStatus=unlisted, got %s",
							body.Status.PrivacyStatus,
						)
					}
					found := false
					for _, tag := range body.Snippet.Tags {
						if tag == "yutu🐰" {
							found = true
						}
					}
					if !found {
						t.Errorf("expected yutu🐰 tag in %v", body.Snippet.Tags)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update video with publishAt on private video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithPublishAt(new("2026-08-18T15:00:00Z")),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "private"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					var body youtube.Video
					if err := json.UnmarshalRead(r.Body, &body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if body.Status.PublishAt != "2026-08-18T15:00:00Z" {
						t.Errorf(
							"expected publishAt=2026-08-18T15:00:00Z, got %s",
							body.Status.PublishAt,
						)
					}
					if body.Status.PrivacyStatus != "private" {
						t.Errorf(
							"expected privacyStatus=private when scheduling, got %s",
							body.Status.PrivacyStatus,
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update video with publishAt and explicit private on public video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithPrivacy("private"),
				WithPublishAt(new("2026-08-18T15:00:00Z")),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					var body youtube.Video
					if err := json.UnmarshalRead(r.Body, &body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if body.Status.PublishAt != "2026-08-18T15:00:00Z" {
						t.Errorf(
							"expected publishAt=2026-08-18T15:00:00Z, got %s",
							body.Status.PublishAt,
						)
					}
					if body.Status.PrivacyStatus != "private" {
						t.Errorf(
							"expected privacyStatus=private when scheduling, got %s",
							body.Status.PrivacyStatus,
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update video with publishAt rejects public without explicit private",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithPublishAt(new("2026-08-18T15:00:00Z")),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					t.Error("expected no update request when scheduling a public video")
				}
			},
			wantErr: true,
		},
		{
			name: "update video with publishAt rejects unlisted without explicit private",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithPublishAt(new("2026-08-18T15:00:00Z")),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "unlisted"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					t.Error("expected no update request when scheduling an unlisted video")
				}
			},
			wantErr: true,
		},
		{
			name: "update video with embeddable, containsSyntheticMedia, recordingDate",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithEmbeddable(&embeddableTrue),
				WithContainsSyntheticMedia(&containsSyntheticMediaTrue),
				WithRecordingDate(new("2024-01-01T00:00:00Z")),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if !strings.Contains(r.URL.Query().Get("part"), "recordingDetails") {
						t.Errorf(
							"expected part to contain recordingDetails, got %s",
							r.URL.Query().Get("part"),
						)
					}
					var body youtube.Video
					if err := json.UnmarshalRead(r.Body, &body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if !body.Status.Embeddable {
						t.Errorf("expected embeddable=true")
					}
					if !body.Status.ContainsSyntheticMedia {
						t.Errorf("expected containsSyntheticMedia=true")
					}
					if body.RecordingDetails == nil || body.RecordingDetails.RecordingDate != "2024-01-01T00:00:00Z" {
						t.Errorf(
							"expected recordingDate=2024-01-01T00:00:00Z, got %+v",
							body.RecordingDetails,
						)
					}
				}
			},
			wantErr: false,
		},
		{
			name: "update video not found",
			opts: []Option{
				WithIds([]string{"missing"}),
				WithTitle("X"),
				WithMaxResults(1),
			},
			getResponse: `{"items": []}`,
			wantErr:     true,
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
								_, _ = w.Write([]byte(tt.getResponse))
							} else {
								_, _ = w.Write([]byte(`{"id": "video-id", "snippet": {"title": "Updated Title"}, "status": {"privacyStatus": "public"}}`))
							}
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_UpdateFalseStatusFields(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status string
		opts   []Option
	}{
		{
			name:   "preserve existing false values",
			status: `{"embeddable":false,"publicStatsViewable":false,"selfDeclaredMadeForKids":false,"containsSyntheticMedia":false}`,
		},
		{
			name:   "explicit false overrides true",
			status: `{"embeddable":true,"publicStatsViewable":false,"selfDeclaredMadeForKids":false,"containsSyntheticMedia":true}`,
			opts:   []Option{WithEmbeddable(new(false)), WithContainsSyntheticMedia(new(false))},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var body struct {
				Status map[string]any `json:"status"`
			}
			svc := common.NewTestService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodGet {
					_, _ = io.WriteString(w, `{"items":[{"id":"video-id","snippet":{"title":"Old title","categoryId":"22"},"status":`+tt.status+`}]}`)
					return
				}
				if err := json.UnmarshalRead(r.Body, &body); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				_, _ = io.WriteString(w, `{"id":"video-id"}`)
			}))
			opts := append([]Option{WithService(svc), WithIds([]string{"video-id"}), WithMaxResults(1), WithTitle("New title"), WithOutput("silent")}, tt.opts...)
			if err := NewVideo(opts...).Update(io.Discard); err != nil {
				t.Fatal(err)
			}
			for _, field := range []string{"embeddable", "publicStatsViewable", "selfDeclaredMadeForKids", "containsSyntheticMedia"} {
				if value, present := body.Status[field]; !present || value != false {
					t.Errorf("status.%s must be present and false, got %v (present=%v)", field, value, present)
				}
			}
		})
	}
}

func TestVideo_Rate(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "rate video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithRating("like"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Query().Get("rating") != "like" {
					t.Errorf("expected rating=like, got %s", r.URL.Query().Get("rating"))
				}
			},
			wantErr: false,
		},
		{
			name: "rate multiple videos",
			opts: []Option{
				WithIds([]string{"video-1", "video-2"}),
				WithRating("dislike"),
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
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Rate(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Rate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_GetRating(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "get rating table",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithOutput("table"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "video-id" {
					t.Errorf("expected id=video-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantErr: false,
		},
		{
			name: "get rating json",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithOutput("json"),
			},
			wantErr: false,
		},
		{
			name: "get rating yaml",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithOutput("yaml"),
			},
			wantErr: false,
		},
		{
			name: "get rating with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithOnBehalfOfContentOwner("owner-id"),
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
						{"videoId": "video-id", "rating": "like"}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.GetRating(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.GetRating() error = %v, wantErr %v", err, tt.wantErr)
				}
				if buf.Len() == 0 {
					t.Errorf("GetRating(%s) output is empty", tt.name)
				}
			},
		)
	}
}

func TestVideo_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete video",
			opts: []Option{
				WithIds([]string{"video-id"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "delete video with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
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
		{
			name: "delete multiple videos",
			opts: []Option{
				WithIds([]string{"video-1", "video-2"}),
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
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVideo_ReportAbuse(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "report abuse",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithReasonId("reason-id"),
				WithComments("spam"),
				WithLanguage(new("en")),
				WithSecondaryReasonId("secondary-reason"),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
			},
			wantErr: false,
		},
		{
			name: "report abuse with onBehalfOfContentOwner",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithReasonId("reason-id"),
				WithOnBehalfOfContentOwner("owner-id"),
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
			},
			wantErr: false,
		},
		{
			name: "report abuse multiple videos",
			opts: []Option{
				WithIds([]string{"video-1", "video-2"}),
				WithReasonId("reason-id"),
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
				v := NewVideo(opts...)
				var buf bytes.Buffer
				if err := v.ReportAbuse(&buf); (err != nil) != tt.wantErr {
					t.Errorf("Video.ReportAbuse() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
