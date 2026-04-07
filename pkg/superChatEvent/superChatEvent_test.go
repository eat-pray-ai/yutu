// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package superChatEvent

import (
	"io"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewSuperChatEvent(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want ISuperChatEvent[youtube.SuperChatEvent]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithHl("en"),
					WithMaxResults(50),
					WithParts([]string{"id", "snippet"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{
					Service:    svc,
					Parts:      []string{"id", "snippet"},
					Output:     "json",
					Hl:         "en",
					MaxResults: 50,
				},
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &SuperChatEvent{Fields: &common.Fields{}},
		},
		{
			name: "with zero max results",
			args: args{
				opts: []Option{
					WithMaxResults(0),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{MaxResults: math.MaxInt64},
			},
		},
		{
			name: "with negative max results",
			args: args{
				opts: []Option{
					WithMaxResults(-10),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{MaxResults: 1},
			},
		},
		{
			name: "with empty string values",
			args: args{
				opts: []Option{
					WithHl(""),
					WithOutput(""),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{Output: "", Hl: ""},
			},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithHl("ja"),
					WithMaxResults(25),
					WithOutput("yaml"),
				},
			},
			want: &SuperChatEvent{
				Fields: &common.Fields{
					Output:     "yaml",
					Hl:         "ja",
					MaxResults: 25,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewSuperChatEvent(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf(
						"%s\nNewSuperChatEvent() = %v\nwant %v", tt.name, got, tt.want,
					)
				}
			},
		)
	}
}

func TestSuperChatEvent_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get super chat events with hl",
			opts: []Option{
				WithHl("ja"),
				WithMaxResults(1),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("hl") != "ja" {
					t.Errorf("expected hl=ja, got %s", r.URL.Query().Get("hl"))
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
						{"id": "event-1", "snippet": {"displayString": "¥500"}}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				s := NewSuperChatEvent(opts...)
				got, err := s.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"SuperChatEvent.Get() error = %v, wantErr %v", err, tt.wantErr,
					)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf(
						"SuperChatEvent.Get() got length = %v, want %v", len(got),
						tt.wantLen,
					)
				}
			},
		)
	}
}

func TestSuperChatEvent_Get_Pagination(t *testing.T) {
	svc := common.NewTestService(t, common.PaginationHandler("event"))
	s := NewSuperChatEvent(WithService(svc), WithMaxResults(22))
	got, err := s.Get()
	if err != nil {
		t.Errorf("SuperChatEvent.Get() error = %v", err)
	}
	if len(got) != 22 {
		t.Errorf("SuperChatEvent.Get() got length = %v, want 22", len(got))
	}
}

func TestSuperChatEvent_List(t *testing.T) {
	common.RunListTest(t, `{
			"items": [
				{
					"id": "event-1",
					"snippet": {
						"displayString": "¥500",
						"commentText": "Hello",
						"supporterDetails": {
							"displayName": "User"
						}
					}
				}
			]
		}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			s := NewSuperChatEvent(WithService(svc), WithOutput(output))
			return s.List
		},
	)
}
