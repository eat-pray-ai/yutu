// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package membershipsLevel

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
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
					WithService(svc),
				},
			},
			want: &MembershipsLevel{
				Fields: &common.Fields{
					Service: svc,
					Parts:   []string{"snippet"},
					Output:  "json",
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
	svc := common.NewTestService(
		t, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(
					[]byte(`{
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
		}`),
				)
			},
		),
	)

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
	common.RunListTest(t, `{
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
		}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			m := NewMembershipsLevel(WithService(svc), WithOutput(output))
			return m.List
		},
	)
}
