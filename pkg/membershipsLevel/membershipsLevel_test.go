// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

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

func TestNewMembershipsLevel(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IMembershipsLevel[youtube.MembershipsLevel]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithJsonpath("items"),
					WithService(svc),
				},
			},
			want: &MembershipsLevel{
				Fields: &common.Fields{
					Service:  svc,
					Parts:    []string{"snippet"},
					Output:   "json",
					Jsonpath: "items",
				},
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &MembershipsLevel{Fields: &common.Fields{}},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewMembershipsLevel(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewMembershipsLevel() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}

func TestMembershipsLevel_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"id": "level-1",
					"snippet": {
						"levelDetails": {
							"displayName": "Level 1"
						}
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

	m := NewMembershipsLevel(WithService(svc))
	got, err := m.Get()
	if err != nil {
		t.Errorf("MembershipsLevel.Get() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("MembershipsLevel.Get() got length = %v, want 1", len(got))
	}
}

func TestMembershipsLevel_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"items": [
				{
					"id": "level-1",
					"snippet": {
						"levelDetails": {
							"displayName": "Level 1"
						}
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
			name: "list memberships levels json",
			opts: []Option{
				WithService(svc),
				WithOutput("json"),
			},
			output:  "json",
			wantErr: false,
		},
		{
			name: "list memberships levels yaml",
			opts: []Option{
				WithService(svc),
				WithOutput("yaml"),
			},
			output:  "yaml",
			wantErr: false,
		},
		{
			name: "list memberships levels table",
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
			m := NewMembershipsLevel(tt.opts...)
			var buf bytes.Buffer
			if err := m.List(&buf); (err != nil) != tt.wantErr {
				t.Errorf("MembershipsLevel.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if buf.Len() == 0 {
				t.Errorf("MembershipsLevel.List() output is empty")
			}
		})
	}
}
