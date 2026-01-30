// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package videoCategory

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

func TestNewVideoCategory(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IVideoCategory[youtube.VideoCategory]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithIds([]string{"cat1", "cat2"}),
					WithHl("en"),
					WithRegionCode("US"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("items.id"),
					WithService(svc),
				},
			},
			want: &VideoCategory{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "items.id",
				},
				Ids:        []string{"cat1", "cat2"},
				Hl:         "en",
				RegionCode: "US",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &VideoCategory{Fields: &common.Fields{}},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
					WithRegionCode(""),
				},
			},
			want: &VideoCategory{
				Fields:     &common.Fields{},
				Hl:         "",
				RegionCode: "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHl("ja"),
					WithRegionCode("JP"),
				},
			},
			want: &VideoCategory{
				Fields:     &common.Fields{},
				Hl:         "ja",
				RegionCode: "JP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewVideoCategory(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewVideoCategory() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestVideoCategory_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get video categories by id",
			opts: []Option{
				WithIds([]string{"category-id"}),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("id") != "category-id" {
					t.Errorf("expected id=category-id, got %s", r.URL.Query().Get("id"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get video categories by regionCode",
			opts: []Option{
				WithRegionCode("US"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("regionCode") != "US" {
					t.Errorf(
						"expected regionCode=US, got %s", r.URL.Query().Get("regionCode"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get video categories with hl",
			opts: []Option{
				WithHl("es"),
				WithIds([]string{"category-id"}),
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
						{"id": "category-1", "snippet": {"title": "Category 1"}}
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
				vc := NewVideoCategory(opts...)
				got, err := vc.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("VideoCategory.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"VideoCategory.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestVideoCategory_List(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
			"items": [
				{
					"id": "category-1",
					"snippet": {
						"title": "Category 1",
						"assignable": true
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
			name: "list video categories json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
				WithIds([]string{"category-1"}),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list video categories yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
				WithIds([]string{"category-1"}),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list video categories table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
				WithIds([]string{"category-1"}),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				vc := NewVideoCategory(tt.opts...)
				var buf bytes.Buffer
				if err := vc.List(&buf); (err != nil) != tt.wantErr {
					t.Errorf(
						"VideoCategory.List() error = %v, wantErr %v", err, tt.wantErr,
					)
				}
				if buf.Len() == 0 {
					t.Errorf("VideoCategory.List() output is empty")
				}
			},
		)
	}
}
