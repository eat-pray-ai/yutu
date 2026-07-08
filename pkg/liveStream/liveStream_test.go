// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveStream

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

func TestNewLiveStream(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
	mine := true

	tests := []struct {
		name string
		args args
		want ILiveStream[youtube.LiveStream]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"stream1", "stream2"}),
					WithTitle("Test Stream"),
					WithDescription("A test stream"),
					WithMine(&mine),
					WithFrameRate("60fps"),
					WithIngestionType("rtmp"),
					WithResolution("1080p"),
					WithMaxResults(50),
					WithParts([]string{"snippet", "cdn", "status"}),
					WithOutput("json"),
					WithService(svc),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("channel123"),
				},
			},
			want: &LiveStream{
				Fields: common.Fields{
					Service:                svc,
					Parts:                  []string{"snippet", "cdn", "status"},
					Output:                 "json",
					MaxResults:             50,
					Ids:                    []string{"stream1", "stream2"},
					OnBehalfOfContentOwner: "owner123",
				},
				Title:                         "Test Stream",
				Description:                   "A test stream",
				Mine:                          &mine,
				FrameRate:                     "60fps",
				IngestionType:                 "rtmp",
				Resolution:                    "1080p",
				OnBehalfOfContentOwnerChannel: "channel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &LiveStream{Fields: common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &LiveStream{
				Fields: common.Fields{MaxResults: math.MaxInt64},
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &LiveStream{
				Fields: common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithTitle(""),
					WithDescription(""),
					WithFrameRate(""),
					WithIngestionType(""),
					WithResolution(""),
				},
			},
			want: &LiveStream{
				Fields:        common.Fields{},
				Title:         "",
				Description:   "",
				FrameRate:     "",
				IngestionType: "",
				Resolution:    "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Stream"),
					WithIngestionType("rtmp"),
					WithMaxResults(25),
					WithParts([]string{"snippet"}),
				},
			},
			want: &LiveStream{
				Fields: common.Fields{
					Parts:      []string{"snippet"},
					MaxResults: 25,
				},
				Title:         "My Stream",
				IngestionType: "rtmp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewLiveStream(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewLiveStream() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestLiveStream_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get live streams by id",
			opts: []Option{
				WithIds([]string{"stream-1"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("id") != "stream-1" {
					t.Errorf(
						"expected id=stream-1, got %s",
						q.Get("id"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get live streams with mine",
			opts: func() []Option {
				mine := true
				return []Option{
					WithMine(&mine),
					WithMaxResults(1),
				}
			}(),
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("mine") != "true" {
					t.Errorf(
						"expected mine=true, got %s",
						q.Get("mine"),
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
						{
							"id": "stream-1",
							"snippet": {
								"title": "Test Stream",
								"description": "A test stream"
							},
							"status": {
								"streamStatus": "active"
							}
						}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewLiveStream(opts...)
				got, err := s.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("LiveStream.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"LiveStream.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestLiveStream_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("stream"))

	mine := true
	s := NewLiveStream(
		WithService(svc),
		WithMine(&mine),
		WithMaxResults(22),
	)
	got, err := s.Get()
	if err != nil {
		t.Errorf("LiveStream.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("LiveStream.Get() got length = %v, want 22", len(got))
	}
}

func TestLiveStream_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "stream-1",
				"snippet": {
					"title": "Test Stream",
					"description": "A test stream"
				},
				"status": {
					"streamStatus": "active"
				}
			}
		]
	}`

	mine := true
	common.RunListTest(
		t, mockResponse,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			s := NewLiveStream(
				WithService(svc),
				WithMine(&mine),
				WithOutput(output),
				WithMaxResults(1),
			)
			return s.List
		},
	)
}

func TestLiveStream_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert live stream",
			opts: []Option{
				WithTitle("New Stream"),
				WithDescription("A new live stream"),
				WithFrameRate("60fps"),
				WithIngestionType("rtmp"),
				WithResolution("1080p"),
				WithParts([]string{"snippet", "cdn"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						Title       string `json:"title"`
						Description string `json:"description"`
					} `json:"snippet"`
					Cdn struct {
						FrameRate     string `json:"frameRate"`
						IngestionType string `json:"ingestionType"`
						Resolution    string `json:"resolution"`
					} `json:"cdn"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Snippet.Title != "New Stream" {
					t.Errorf(
						"expected snippet.title=New Stream, got %s", body.Snippet.Title,
					)
				}
				if body.Snippet.Description != "A new live stream" {
					t.Errorf(
						"expected snippet.description=A new live stream, got %s",
						body.Snippet.Description,
					)
				}
				if body.Cdn.FrameRate != "60fps" {
					t.Errorf("expected cdn.frameRate=60fps, got %s", body.Cdn.FrameRate)
				}
				if body.Cdn.IngestionType != "rtmp" {
					t.Errorf(
						"expected cdn.ingestionType=rtmp, got %s", body.Cdn.IngestionType,
					)
				}
				if body.Cdn.Resolution != "1080p" {
					t.Errorf("expected cdn.resolution=1080p, got %s", body.Cdn.Resolution)
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
							_, _ = w.Write([]byte(`{"id": "new-stream-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewLiveStream(opts...)
				var buf bytes.Buffer
				if err := s.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveStream.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveStream_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "update live stream",
			opts: []Option{
				WithIds([]string{"stream-1"}),
				WithTitle("Updated Stream"),
				WithDescription("Updated description"),
				WithMaxResults(1),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				callCount := 0
				svc := common.NewTestService(
					t, http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json")
							callCount++
							if callCount == 1 {
								_, _ = w.Write(
									[]byte(`{
									"items": [{
										"id": "stream-1",
										"snippet": {"title": "Old Stream", "description": "Old description"},
										"cdn": {"frameRate": "30fps", "ingestionType": "rtmp", "resolution": "720p"}
									}]
								}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "stream-1"}`))
							}
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewLiveStream(opts...)
				var buf bytes.Buffer
				if err := s.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveStream.Update() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveStream_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete live stream",
			opts: []Option{
				WithIds([]string{"stream-id"}),
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
				s := NewLiveStream(opts...)
				var buf bytes.Buffer
				if err := s.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveStream.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
