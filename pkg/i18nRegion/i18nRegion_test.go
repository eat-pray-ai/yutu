// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package i18nRegion

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewI18nRegion(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want II18nRegion[youtube.I18nRegion]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &I18nRegion{
				Fields: &common.Fields{
					Service: svc,
					Parts:   []string{"id", "snippet"},
					Output:  "json",
					Hl:      "en",
				},
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &I18nRegion{Fields: &common.Fields{}},
		},
		{
			name: "with empty string value",
			args: args{
				opts: []Option{
					WithHl(""),
				},
			},
			want: &I18nRegion{
				Fields: &common.Fields{Hl: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewI18nRegion(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewI18nRegion() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestI18nRegion_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get i18n regions with hl",
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
								[]byte(`{"items": [{"id": "US", "snippet": {"gl": "US", "name": "United States"}}]}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				i := NewI18nRegion(opts...)
				got, err := i.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("I18nRegion.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"I18nRegion.Get() got length = %v, want %v", len(got), tt.wantLen,
					)
				}
			},
		)
	}
}

func TestI18nRegion_List(t *testing.T) {
	common.RunListTest(t, `{
			"items": [
				{
					"id": "US",
					"snippet": {
						"gl": "US",
						"name": "United States"
					}
				}
			]
		}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			i := NewI18nRegion(WithService(svc), WithOutput(output))
			return i.List
		},
	)
}
