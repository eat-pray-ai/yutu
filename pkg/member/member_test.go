// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package member

import (
	"io"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewMember(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IMember[youtube.Member]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithMemberChannelId("member123"),
					WithHasAccessToLevel("level1"),
					WithMaxResults(100),
					WithMode("all_current"),
					WithParts([]string{"snippet"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &Member{
				Fields: &common.Fields{
					Service:    svc,
					Parts:      []string{"snippet"},
					Output:     "json",
					MaxResults: 100,
				},
				MemberChannelId:  "member123",
				HasAccessToLevel: "level1",
				Mode:             "all_current",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &Member{Fields: &common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &Member{
				Fields: &common.Fields{MaxResults: math.MaxInt64},
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-15),
				},
			},
			want: &Member{
				Fields: &common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithMemberChannelId(""),
					WithHasAccessToLevel(""),
					WithMode(""),
				},
			},
			want: &Member{
				Fields:           &common.Fields{},
				MemberChannelId:  "",
				HasAccessToLevel: "",
				Mode:             "",
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithMemberChannelId("channel456"),
					WithMaxResults(50),
					WithParts([]string{"id"}),
				},
			},
			want: &Member{
				Fields: &common.Fields{
					Parts:      []string{"id"},
					MaxResults: 50,
				},
				MemberChannelId: "channel456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewMember(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("%s\nNewMember() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestMember_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get members with memberChannelId",
			opts: []Option{
				WithMemberChannelId("channel-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				q := r.URL.Query()
				if q.Get("filterByMemberChannelId") != "channel-id" {
					t.Errorf(
						"expected filterByMemberChannelId=channel-id, got %s",
						q.Get("filterByMemberChannelId"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get members with hasAccessToLevel",
			opts: []Option{
				WithHasAccessToLevel("level-id"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hasAccessToLevel") != "level-id" {
					t.Errorf(
						"expected hasAccessToLevel=level-id, got %s",
						r.URL.Query().Get("hasAccessToLevel"),
					)
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get members with mode",
			opts: []Option{
				WithMode("all_current"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("mode") != "all_current" {
					t.Errorf(
						"expected mode=all_current, got %s", r.URL.Query().Get("mode"),
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
							"snippet": {
								"memberDetails": {
									"channelId": "channel-1",
									"displayName": "Member 1"
								}
							}
						}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				m := NewMember(opts...)
				got, err := m.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("Member.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("Member.Get() got length = %v, want %v", len(got), tt.wantLen)
				}
			},
		)
	}
}

func TestMember_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("member"))

	m := NewMember(
		WithService(svc),
		WithMaxResults(22),
	)
	got, err := m.Get()
	if err != nil {
		t.Errorf("Member.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("Member.Get() got length = %v, want 22", len(got))
	}
}

func TestMember_List(t *testing.T) {
	mockResponse := `{
		"items": [
			{
				"snippet": {
					"memberDetails": {
						"channelId": "channel-1",
						"displayName": "Member 1"
					}
				}
			}
		]
	}`

	common.RunListTest(
		t, mockResponse,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			m := NewMember(
				WithService(svc),
				WithOutput(output),
				WithMaxResults(1),
			)
			return m.List
		},
	)
}
