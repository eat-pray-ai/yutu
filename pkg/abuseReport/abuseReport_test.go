// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package abuseReport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewAbuseReport(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IAbuseReport
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithAbuseTypes([]string{"spam", "harassment"}),
					WithDescription("This video contains spam"),
					WithSubjectId("video123"),
					WithSubjectTypeId("video"),
					WithSubjectUrl("https://youtube.com/watch?v=video123"),
					WithRelatedEntityId("channel456"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &AbuseReport{
				Fields: common.Fields{
					Service: svc,
					Parts:   []string{"snippet"},
					Output:  "json",
				},
				AbuseTypes:      []string{"spam", "harassment"},
				Description:     "This video contains spam",
				SubjectId:       "video123",
				SubjectTypeId:   "video",
				SubjectUrl:      "https://youtube.com/watch?v=video123",
				RelatedEntityId: "channel456",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &AbuseReport{Fields: common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithDescription(""),
					WithSubjectId(""),
					WithSubjectTypeId(""),
				},
			},
			want: &AbuseReport{
				Fields:        common.Fields{},
				Description:   "",
				SubjectId:     "",
				SubjectTypeId: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithAbuseTypes([]string{"violence"}),
					WithSubjectId("video789"),
					WithDescription("Violent content"),
				},
			},
			want: &AbuseReport{
				Fields:      common.Fields{},
				AbuseTypes:  []string{"violence"},
				SubjectId:   "video789",
				Description: "Violent content",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewAbuseReport(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewAbuseReport() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestAbuseReport_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert abuse report",
			opts: []Option{
				WithAbuseTypes([]string{"spam"}),
				WithDescription("This is spam"),
				WithSubjectId("video123"),
				WithSubjectTypeId("video"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					AbuseTypes []struct {
						Id string `json:"id"`
					} `json:"abuseTypes"`
					Description string `json:"description"`
					Subject     struct {
						Id     string `json:"id"`
						TypeId string `json:"typeId"`
					} `json:"subject"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if len(body.AbuseTypes) != 1 || body.AbuseTypes[0].Id != "spam" {
					t.Errorf("expected abuseTypes=[{id:spam}], got %v", body.AbuseTypes)
				}
				if body.Description != "This is spam" {
					t.Errorf("expected description=This is spam, got %s", body.Description)
				}
				if body.Subject.Id != "video123" {
					t.Errorf("expected subject.id=video123, got %s", body.Subject.Id)
				}
				if body.Subject.TypeId != "video" {
					t.Errorf("expected subject.typeId=video, got %s", body.Subject.TypeId)
				}
			},
			wantErr: false,
		},
		{
			name: "insert abuse report with related entity",
			opts: []Option{
				WithAbuseTypes([]string{"harassment"}),
				WithSubjectId("comment456"),
				WithSubjectTypeId("comment"),
				WithRelatedEntityId("channel789"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				defer func() { _ = r.Body.Close() }()
				var body struct {
					RelatedEntities []struct {
						Entity struct {
							Id string `json:"id"`
						} `json:"entity"`
					} `json:"relatedEntities"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if len(body.RelatedEntities) != 1 || body.RelatedEntities[0].Entity.Id != "channel789" {
					t.Errorf(
						"expected relatedEntities=[{entity:{id:channel789}}], got %v",
						body.RelatedEntities,
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
							_, _ = w.Write([]byte(`{}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				a := NewAbuseReport(opts...)
				var buf bytes.Buffer
				if err := a.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"AbuseReport.Insert() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
			},
		)
	}
}
