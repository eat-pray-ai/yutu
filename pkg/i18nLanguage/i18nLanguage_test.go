// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nLanguage

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

func TestNewI18nLanguage(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want II18nLanguage[youtube.I18nLanguage]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("$"),
					WithService(svc),
				},
			},
			want: &I18nLanguage{
				Fields: &common.Fields{
					Service:  &youtube.Service{},
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "$",
				},
				Hl: "en",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &I18nLanguage{Fields: &common.Fields{}},
		},
		{
			name: "with empty string value",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &I18nLanguage{
				Fields: &common.Fields{},
				Hl:     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewI18nLanguage(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewI18nLanguage() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestI18nLanguage_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get i18n languages with hl",
			opts: []Option{
				WithHl("es_ES"),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "es_ES" {
					t.Errorf("expected hl=es_ES, got %s", r.URL.Query().Get("hl"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.verify != nil {
					tt.verify(r)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{
					"items": [
						{"id": "en", "snippet": {"hl": "en", "name": "English"}}
					]
				}`))
			}))
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
			i := NewI18nLanguage(opts...)
			got, err := i.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nLanguage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("I18nLanguage.Get() got length = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestI18nLanguage_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"id": "en",
					"snippet": {
						"hl": "en",
						"name": "English"
					}
				}
			]
		}`))
	}))
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
			name: "list i18n languages json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list i18n languages yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list i18n languages table",
			opts: []Option{
				WithService(svc),
				WithOutput("table"),
			},
			output:  "table",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewI18nLanguage(tt.opts...)
			var buf bytes.Buffer
			if err := i.List(&buf); (err != nil) != tt.wantErr {
				t.Errorf("I18nLanguage.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if buf.Len() == 0 {
				t.Errorf("I18nLanguage.List() output is empty")
			}
		})
	}
}
