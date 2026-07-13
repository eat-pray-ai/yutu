// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package liveBroadcast

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

func TestNewLiveBroadcast(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}
	mine := true

	tests := []struct {
		name string
		args args
		want ILiveBroadcast[youtube.LiveBroadcast]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"broadcast1", "broadcast2"}),
					WithTitle("Test Broadcast"),
					WithDescription("A test broadcast"),
					WithMine(&mine),
					WithBroadcastStatus("active"),
					WithBroadcastType("event"),
					WithPrivacyStatus("public"),
					WithScheduledStartTime("2026-01-01T00:00:00Z"),
					WithScheduledEndTime("2026-01-01T01:00:00Z"),
					WithStreamId("stream123"),
					WithCueType("cueTypeAd"),
					WithCueDurationSecs(30),
					WithCueInsertionOffsetMs(5000),
					WithCueWalltimeMs(1234567890),
					WithMaxResults(50),
					WithParts([]string{"snippet", "status", "contentDetails"}),
					WithOutput("json"),
					WithService(svc),
					WithOnBehalfOfContentOwner("owner123"),
					WithOnBehalfOfContentOwnerChannel("channel123"),
				},
			},
			want: &LiveBroadcast{
				Fields: common.Fields{
					Service:                svc,
					Parts:                  []string{"snippet", "status", "contentDetails"},
					Output:                 "json",
					MaxResults:             50,
					Ids:                    []string{"broadcast1", "broadcast2"},
					OnBehalfOfContentOwner: "owner123",
				},
				Title:                         "Test Broadcast",
				Description:                   "A test broadcast",
				Mine:                          &mine,
				BroadcastStatus:               "active",
				BroadcastType:                 "event",
				PrivacyStatus:                 "public",
				ScheduledStartTime:            "2026-01-01T00:00:00Z",
				ScheduledEndTime:              "2026-01-01T01:00:00Z",
				StreamId:                      "stream123",
				CueType:                       "cueTypeAd",
				CueDurationSecs:               30,
				CueInsertionOffsetMs:          5000,
				CueWalltimeMs:                 1234567890,
				OnBehalfOfContentOwnerChannel: "channel123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &LiveBroadcast{Fields: common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &LiveBroadcast{
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
			want: &LiveBroadcast{
				Fields: common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithTitle(""),
					WithDescription(""),
					WithBroadcastStatus(""),
					WithBroadcastType(""),
					WithPrivacyStatus(""),
					WithScheduledStartTime(""),
					WithScheduledEndTime(""),
					WithStreamId(""),
				},
			},
			want: &LiveBroadcast{
				Fields:             common.Fields{},
				Title:              "",
				Description:        "",
				BroadcastStatus:    "",
				BroadcastType:      "",
				PrivacyStatus:      "",
				ScheduledStartTime: "",
				ScheduledEndTime:   "",
				StreamId:           "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithTitle("My Broadcast"),
					WithPrivacyStatus("private"),
					WithMaxResults(25),
					WithParts([]string{"snippet"}),
				},
			},
			want: &LiveBroadcast{
				Fields: common.Fields{
					Parts:      []string{"snippet"},
					MaxResults: 25,
				},
				Title:         "My Broadcast",
				PrivacyStatus: "private",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewLiveBroadcast(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewLiveBroadcast() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestLiveBroadcast_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get live broadcasts by id",
			opts: []Option{
				WithIds([]string{"broadcast-1"}),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("id") != "broadcast-1" {
					t.Errorf(
						"expected id=broadcast-1, got %s",
						q.Get("id"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get live broadcasts with mine",
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
		{
			name: "get live broadcasts with status filter",
			opts: []Option{
				WithBroadcastStatus("active"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("broadcastStatus") != "active" {
					t.Errorf(
						"expected broadcastStatus=active, got %s",
						q.Get("broadcastStatus"),
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
							"id": "broadcast-1",
							"snippet": {
								"title": "Test Broadcast",
								"description": "A test broadcast"
							},
							"status": {
								"lifeCycleStatus": "ready",
								"privacyStatus": "public"
							}
						}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveBroadcast(opts...)
				got, err := b.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("LiveBroadcast.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"LiveBroadcast.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestLiveBroadcast_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("broadcast"))

	mine := true
	b := NewLiveBroadcast(
		WithService(svc),
		WithMine(&mine),
		WithMaxResults(22),
	)
	got, err := b.Get()
	if err != nil {
		t.Errorf("LiveBroadcast.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("LiveBroadcast.Get() got length = %v, want 22", len(got))
	}
}

func TestLiveBroadcast_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"id": "broadcast-1",
				"snippet": {
					"title": "Test Broadcast",
					"description": "A test broadcast"
				},
				"status": {
					"lifeCycleStatus": "ready",
					"privacyStatus": "public"
				}
			}
		]
	}`

	mine := true
	common.RunListTest(
		t, mockResponse,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			b := NewLiveBroadcast(
				WithService(svc),
				WithMine(&mine),
				WithOutput(output),
				WithMaxResults(1),
			)
			return b.List
		},
	)
}

func TestLiveBroadcast_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert live broadcast",
			opts: []Option{
				WithTitle("New Broadcast"),
				WithDescription("A new live broadcast"),
				WithPrivacyStatus("public"),
				WithScheduledStartTime("2026-01-01T00:00:00Z"),
				WithScheduledEndTime("2026-01-01T01:00:00Z"),
				WithParts([]string{"snippet", "status"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					Snippet struct {
						Title              string `json:"title"`
						Description        string `json:"description"`
						ScheduledStartTime string `json:"scheduledStartTime"`
						ScheduledEndTime   string `json:"scheduledEndTime"`
					} `json:"snippet"`
					Status struct {
						PrivacyStatus string `json:"privacyStatus"`
					} `json:"status"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.Snippet.Title != "New Broadcast" {
					t.Errorf(
						"expected snippet.title=New Broadcast, got %s", body.Snippet.Title,
					)
				}
				if body.Snippet.Description != "A new live broadcast" {
					t.Errorf(
						"expected snippet.description=A new live broadcast, got %s",
						body.Snippet.Description,
					)
				}
				if body.Status.PrivacyStatus != "public" {
					t.Errorf(
						"expected status.privacyStatus=public, got %s",
						body.Status.PrivacyStatus,
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
							_, _ = w.Write([]byte(`{"id": "new-broadcast-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveBroadcast(opts...)
				var buf bytes.Buffer
				if err := b.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveBroadcast.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveBroadcast_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name: "update live broadcast",
			opts: []Option{
				WithIds([]string{"broadcast-1"}),
				WithTitle("Updated Broadcast"),
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
										"id": "broadcast-1",
										"snippet": {
											"title": "Old Broadcast",
											"description": "Old description",
											"scheduledStartTime": "2026-01-01T00:00:00Z",
											"scheduledEndTime": "2026-01-01T01:00:00Z"
										},
										"status": {"privacyStatus": "public", "lifeCycleStatus": "ready"}
									}]
								}`),
								)
							} else {
								_, _ = w.Write([]byte(`{"id": "broadcast-1"}`))
							}
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveBroadcast(opts...)
				var buf bytes.Buffer
				if err := b.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveBroadcast.Update() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveBroadcast_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete live broadcast",
			opts: []Option{
				WithIds([]string{"broadcast-id"}),
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
				b := NewLiveBroadcast(opts...)
				var buf bytes.Buffer
				if err := b.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveBroadcast.Delete() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveBroadcast_Bind(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "bind live broadcast to stream",
			opts: []Option{
				WithIds([]string{"broadcast-id"}),
				WithStreamId("stream-id"),
				WithParts([]string{"id", "snippet", "contentDetails", "status"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				q := r.URL.Query()
				if q.Get("id") != "broadcast-id" {
					t.Errorf("expected id=broadcast-id, got %s", q.Get("id"))
				}
				if q.Get("streamId") != "stream-id" {
					t.Errorf("expected streamId=stream-id, got %s", q.Get("streamId"))
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
							_, _ = w.Write([]byte(`{"id": "broadcast-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveBroadcast(opts...)
				var buf bytes.Buffer
				if err := b.Bind(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveBroadcast.Bind() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveBroadcast_Transition(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "transition live broadcast",
			opts: []Option{
				WithIds([]string{"broadcast-id"}),
				WithBroadcastStatus("live"),
				WithParts([]string{"id", "snippet", "status"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				q := r.URL.Query()
				if q.Get("id") != "broadcast-id" {
					t.Errorf("expected id=broadcast-id, got %s", q.Get("id"))
				}
				if q.Get("broadcastStatus") != "live" {
					t.Errorf(
						"expected broadcastStatus=live, got %s",
						q.Get("broadcastStatus"),
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
							_, _ = w.Write([]byte(`{"id": "broadcast-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveBroadcast(opts...)
				var buf bytes.Buffer
				if err := b.Transition(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveBroadcast.Transition() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}

func TestLiveBroadcast_InsertCuepoint(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert cuepoint",
			opts: []Option{
				WithIds([]string{"broadcast-id"}),
				WithCueType("cueTypeAd"),
				WithCueDurationSecs(30),
				WithCueInsertionOffsetMs(5000),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}
				q := r.URL.Query()
				if q.Get("id") != "broadcast-id" {
					t.Errorf("expected id=broadcast-id, got %s", q.Get("id"))
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					CueType      string `json:"cueType"`
					DurationSecs int64  `json:"durationSecs"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}
				if body.CueType != "cueTypeAd" {
					t.Errorf("expected cueType=cueTypeAd, got %s", body.CueType)
				}
				if body.DurationSecs != 30 {
					t.Errorf("expected durationSecs=30, got %d", body.DurationSecs)
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
							_, _ = w.Write([]byte(`{"id": "cuepoint-id"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				b := NewLiveBroadcast(opts...)
				var buf bytes.Buffer
				if err := b.InsertCuepoint(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"LiveBroadcast.InsertCuepoint() error = %v, wantErr %v", err,
						tt.wantErr,
					)
				}
			},
		)
	}
}
