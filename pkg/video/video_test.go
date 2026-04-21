// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package video

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewVideo(t *testing.T) {
	type args struct {
		opts []Option
	}

	autoLevelsTrue := true
	autoLevelsFalse := false
	forKidsTrue := true
	forKidsFalse := false
	embeddableTrue := true
	embeddableFalse := false
	containsSyntheticMediaTrue := true
	containsSyntheticMediaFalse := false
	stabilizeTrue := true
	stabilizeFalse := false
	notifySubscribersTrue := true
	notifySubscribersFalse := false
	publicStatsViewableTrue := true
	publicStatsViewableFalse := false
	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IVideo[youtube.Video]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"video1", "video2"}),
					WithAutoLevels(&autoLevelsTrue),
					WithFile("/path/to/video.mp4"),
					WithTitle("Test Video"),
					WithDescription("Test video description"),
					WithHl("en"),
					WithTags([]string{"tag1", "tag2", "tag3"}),
					WithLanguage("en"),
					WithLocale("en_US"),
					WithLicense("youtube"),
					WithThumbnail("/path/to/thumbnail.jpg"),
					WithRating("like"),
					WithChart("mostPopular"),
					WithChannelId("channel123"),
					WithComments("Test comments"),
					WithPlaylistId("playlist123"),
					WithCategory("22"),
					WithPrivacy("public"),
					WithForKids(&forKidsTrue),
					WithEmbeddable(&embeddableTrue),
					WithContainsSyntheticMedia(&containsSyntheticMediaTrue),
					WithRecordingDate("2024-06-15T10:00:00Z"),
					WithPublishAt("2024-12-31T23:59:59Z"),
					WithRegionCode("US"),
					WithReasonId("reason123"),
					WithSecondaryReasonId("secondaryReason123"),
					WithStabilize(&stabilizeTrue),
					WithMaxHeight(1080),
					WithMaxWidth(1920),
					WithMaxResults(50),
					WithNotifySubscribers(&notifySubscribersTrue),
					WithPublicStatsViewable(&publicStatsViewableTrue),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("ownerChannel123"),
					WithParts([]string{"snippet", "contentDetails"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &Video{
				Fields: &common.Fields{
					Service:                svc,
					Parts:                  []string{"snippet", "contentDetails"},
					Output:                 "json",
					Ids:                    []string{"video1", "video2"},
					MaxResults:             50,
					Hl:                     "en",
					ChannelId:              "channel123",
					OnBehalfOfContentOwner: "owner123",
				},
				AutoLevels:                    &autoLevelsTrue,
				File:                          "/path/to/video.mp4",
				Title:                         "Test Video",
				Description:                   "Test video description",
				Tags:                          []string{"tag1", "tag2", "tag3"},
				Language:                      "en",
				Locale:                        "en_US",
				License:                       "youtube",
				Thumbnail:                     "/path/to/thumbnail.jpg",
				Rating:                        "like",
				Chart:                         "mostPopular",
				Comments:                      "Test comments",
				PlaylistId:                    "playlist123",
				CategoryId:                    "22",
				Privacy:                       "public",
				ForKids:                       &forKidsTrue,
				Embeddable:                    &embeddableTrue,
				ContainsSyntheticMedia:        &containsSyntheticMediaTrue,
				RecordingDate:                 "2024-06-15T10:00:00Z",
				PublishAt:                     "2024-12-31T23:59:59Z",
				RegionCode:                    "US",
				ReasonId:                      "reason123",
				SecondaryReasonId:             "secondaryReason123",
				Stabilize:                     &stabilizeTrue,
				MaxHeight:                     1080,
				MaxWidth:                      1920,
				NotifySubscribers:             &notifySubscribersTrue,
				PublicStatsViewable:           &publicStatsViewableTrue,
				OnBehalfOfContentOwnerChannel: "ownerChannel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Video{Fields: &common.Fields{}},
		},
		{
			name: "with nil boolean options",
			args: args{
				opts: []Option{
					WithAutoLevels(nil),
					WithForKids(nil),
					WithEmbeddable(nil),
					WithContainsSyntheticMedia(nil),
					WithStabilize(nil),
					WithNotifySubscribers(nil),
					WithPublicStatsViewable(nil),
				},
			},
			want: &Video{Fields: &common.Fields{}},
		},
		{
			name: "with false boolean options",
			args: args{
				opts: []Option{
					WithAutoLevels(&autoLevelsFalse),
					WithForKids(&forKidsFalse),
					WithEmbeddable(&embeddableFalse),
					WithContainsSyntheticMedia(&containsSyntheticMediaFalse),
					WithStabilize(&stabilizeFalse),
					WithNotifySubscribers(&notifySubscribersFalse),
					WithPublicStatsViewable(&publicStatsViewableFalse),
				},
			},
			want: &Video{
				Fields:                 &common.Fields{},
				AutoLevels:             &autoLevelsFalse,
				ForKids:                &forKidsFalse,
				Embeddable:             &embeddableFalse,
				ContainsSyntheticMedia: &containsSyntheticMediaFalse,
				Stabilize:              &stabilizeFalse,
				NotifySubscribers:      &notifySubscribersFalse,
				PublicStatsViewable:    &publicStatsViewableFalse,
			},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Video{
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
			want: &Video{
				Fields: &common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithFile(""),
					WithTitle(""),
					WithDescription(""),
					WithHl(""),
					WithLanguage(""),
					WithLocale(""),
					WithLicense(""),
					WithThumbnail(""),
					WithRating(""),
					WithChart(""),
					WithChannelId(""),
					WithComments(""),
					WithPlaylistId(""),
					WithCategory(""),
					WithPrivacy(""),
					WithRecordingDate(""),
					WithPublishAt(""),
					WithRegionCode(""),
					WithReasonId(""),
					WithSecondaryReasonId(""),
					WithOnBehalfOfContentOwner(""),
					WithOnBehalfOfContentOwnerChannel(""),
				},
			},
			want: &Video{
				Fields:                        &common.Fields{},
				File:                          "",
				Title:                         "",
				Description:                   "",
				Language:                      "",
				Locale:                        "",
				License:                       "",
				Thumbnail:                     "",
				Rating:                        "",
				Chart:                         "",
				Comments:                      "",
				PlaylistId:                    "",
				CategoryId:                    "",
				Privacy:                       "",
				RecordingDate:                 "",
				PublishAt:                     "",
				RegionCode:                    "",
				ReasonId:                      "",
				SecondaryReasonId:             "",
				OnBehalfOfContentOwnerChannel: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Video"),
					WithDescription("A great video"),
					WithTags([]string{"tutorial", "golang"}),
					WithPrivacy("private"),
					WithMaxResults(25),
					WithForKids(&forKidsFalse),
				},
			},
			want: &Video{
				Fields:      &common.Fields{MaxResults: 25},
				Title:       "My Video",
				Description: "A great video",
				Tags:        []string{"tutorial", "golang"},
				Privacy:     "private",
				ForKids:     &forKidsFalse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideo(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewVideo() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
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
		if err := json.NewDecoder(part).Decode(&body); err != nil {
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
				WithRecordingDate("2024-06-15T10:00:00Z"),
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
					t.Errorf("expected part to contain recordingDetails, got %s", r.URL.Query().Get("part"))
				}
				body := decodeMultipartVideo(t, r)
				if body.RecordingDetails == nil || body.RecordingDetails.RecordingDate != "2024-06-15T10:00:00Z" {
					t.Errorf("expected recordingDetails.recordingDate=2024-06-15T10:00:00Z, got %+v", body.RecordingDetails)
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
					t.Errorf("expected snippet.title=test_video (from filename), got %+v", body.Snippet)
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
					t.Errorf("expected exactly 1 yutu🐰 tag, got %d in %v", count, body.Snippet.Tags)
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

	err = os.WriteFile(tmpDir+"/test_video.mp4", []byte("dummy video content"), 0644)
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
		name         string
		opts         []Option
		getResponse  string
		verify       func(*http.Request)
		wantErr      bool
	}{
		{
			name: "update video",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithTitle("Updated Title"),
				WithDescription("Updated Description"),
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
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
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
				WithLanguage("ja"),
				WithLicense("creativeCommon"),
				WithCategory("22"),
				WithPrivacy("unlisted"),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public", "license": "youtube", "embeddable": true, "publicStatsViewable": true}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					var body youtube.Video
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if body.Snippet.DefaultLanguage != "ja" {
						t.Errorf("expected defaultLanguage=ja, got %s", body.Snippet.DefaultLanguage)
					}
					if body.Status.License != "creativeCommon" {
						t.Errorf("expected license=creativeCommon, got %s", body.Status.License)
					}
					if body.Snippet.CategoryId != "22" {
						t.Errorf("expected categoryId=22, got %s", body.Snippet.CategoryId)
					}
					if body.Status.PrivacyStatus != "unlisted" {
						t.Errorf("expected privacyStatus=unlisted, got %s", body.Status.PrivacyStatus)
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
			name: "update video with embeddable, containsSyntheticMedia, recordingDate",
			opts: []Option{
				WithIds([]string{"video-id"}),
				WithEmbeddable(&embeddableTrue),
				WithContainsSyntheticMedia(&containsSyntheticMediaTrue),
				WithRecordingDate("2024-01-01T00:00:00Z"),
				WithMaxResults(1),
			},
			getResponse: `{"items": [{"id": "video-id", "snippet": {"title": "Old Title"}, "status": {"privacyStatus": "public"}}]}`,
			verify: func(r *http.Request) {
				if r.Method == "PUT" {
					if !strings.Contains(r.URL.Query().Get("part"), "recordingDetails") {
						t.Errorf("expected part to contain recordingDetails, got %s", r.URL.Query().Get("part"))
					}
					var body youtube.Video
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Fatalf("failed to decode update body: %v", err)
					}
					if !body.Status.Embeddable {
						t.Errorf("expected embeddable=true")
					}
					if !body.Status.ContainsSyntheticMedia {
						t.Errorf("expected containsSyntheticMedia=true")
					}
					if body.RecordingDetails == nil || body.RecordingDetails.RecordingDate != "2024-01-01T00:00:00Z" {
						t.Errorf("expected recordingDate=2024-01-01T00:00:00Z, got %+v", body.RecordingDetails)
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
				WithLanguage("en"),
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
