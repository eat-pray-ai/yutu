// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoAbuseReportReason

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func TestNewVideoAbuseReportReason(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IVideoAbuseReportReason[youtube.VideoAbuseReportReason]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHL("en"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithJsonpath("$.items[*].id"),
					WithService(svc),
				},
			},
			want: &VideoAbuseReportReason{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"id", "snippet"},
					Output:   "json",
					Jsonpath: "$.items[*].id",
				},
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &VideoAbuseReportReason{Fields: &common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHL(""),
					WithOutput(""),
					WithJsonpath(""),
				},
			},
			want: &VideoAbuseReportReason{
				Fields: &common.Fields{
					Output: "", Jsonpath: "",
				},
				Hl: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHL("ja"),
					WithParts([]string{"snippet"}),
				},
			},
			want: &VideoAbuseReportReason{
				Fields: &common.Fields{Parts: []string{"snippet"}},
				Hl:     "ja",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoAbuseReportReason(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewVideoAbuseReportReason() = %v\nwant %v",
						tt.name, got, tt.want,
					)
				}
			},
		)
	}
}

func TestVideoAbuseReportReason_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get video abuse report reasons",
			opts: []Option{
				WithHL("es"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "es" {
					t.Errorf("expected hl=es, got %s", r.URL.Query().Get("hl"))
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
						{"id": "reason-1", "snippet": {"label": "Reason 1"}}
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
				va := NewVideoAbuseReportReason(opts...)
				got, err := va.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"VideoAbuseReportReason.Get() error = %v, wantErr %v", err,
						tt.wantErr,
					)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"VideoAbuseReportReason.Get() got length = %v, want %v", len(got),
						tt.wantLen,
					)
				}
			},
		)
	}
}

func TestVideoAbuseReportReason_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "reason-1",
					"snippet": {
						"label": "Reason 1"
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
			name: "list video abuse report reasons json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list video abuse report reasons yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list video abuse report reasons table",
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
				va := NewVideoAbuseReportReason(tt.opts...)
				var buf bytes.Buffer
				if err := va.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"VideoAbuseReportReason.List() error = %v, wantErr %v", err,
						tt.wantErr,
					)
				}
				if buf.Len() == 0 {
					t.Errorf("VideoAbuseReportReason.List() output is empty")
				}
			},
		)
	}
}
