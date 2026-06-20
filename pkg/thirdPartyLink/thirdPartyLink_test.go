// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package thirdPartyLink

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/eat-pray-ai/yutu/pkg/common"
	"google.golang.org/api/youtube/v3"
)

func TestNewThirdPartyLink(t *testing.T) {
	type args struct {
		opts []Option
	}

	svc := &youtube.Service{}

	tests := []struct {
		name string
		args args
		want IThirdPartyLink[youtube.ThirdPartyLink]
	}{
		{
			name: "with all options",
			args: args{
				opts: []Option{
					WithLinkingToken("token123"),
					WithType("channelToStoreLink"),
					WithLinkStatus("linked"),
					WithExternalChannelId("ext-channel-123"),
					WithParts([]string{"snippet", "status"}),
					WithOutput("json"),
					WithService(svc),
				},
			},
			want: &ThirdPartyLink{
				Fields: common.Fields{
					Service: svc,
					Parts:   []string{"snippet", "status"},
					Output:  "json",
				},
				LinkingToken:      "token123",
				Type:              "channelToStoreLink",
				LinkStatus:        "linked",
				ExternalChannelId: "ext-channel-123",
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &ThirdPartyLink{Fields: common.Fields{}},
		},
		{
			name: "with partial options",
			args: args{
				opts: []Option{
					WithLinkingToken("token456"),
					WithType("channelToStoreLink"),
				},
			},
			want: &ThirdPartyLink{
				Fields:       common.Fields{},
				LinkingToken: "token456",
				Type:         "channelToStoreLink",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := NewThirdPartyLink(tt.args.opts...); !reflect.DeepEqual(
					got, tt.want,
				) {
					t.Errorf("%s\nNewThirdPartyLink() = %v\nwant %v", tt.name, got, tt.want)
				}
			},
		)
	}
}

func TestThirdPartyLink_Get(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantLen int
		wantErr bool
	}{
		{
			name: "get with linking token",
			opts: []Option{
				WithLinkingToken("token123"),
				WithParts([]string{"snippet", "status"}),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("linkingToken") != "token123" {
					t.Errorf("expected linkingToken=token123, got %s", r.URL.Query().Get("linkingToken"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get with type filter",
			opts: []Option{
				WithType("channelToStoreLink"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("type") != "channelToStoreLink" {
					t.Errorf("expected type=channelToStoreLink, got %s", r.URL.Query().Get("type"))
				}
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "get with external channel id",
			opts: []Option{
				WithExternalChannelId("ext-channel-123"),
				WithParts([]string{"snippet"}),
			},
			verify: func(r *http.Request) {
				if r.URL.Query().Get("externalChannelId") != "ext-channel-123" {
					t.Errorf("expected externalChannelId=ext-channel-123, got %s", r.URL.Query().Get("externalChannelId"))
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
						{"linkingToken": "token123", "snippet": {"type": "channelToStoreLink"}, "status": {"linkStatus": "linked"}}
					]
				}`),
							)
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				tpl := NewThirdPartyLink(opts...)
				got, err := tpl.Get()
				if (err != nil) != tt.wantErr {
					t.Errorf("ThirdPartyLink.Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(got) != tt.wantLen {
					t.Errorf("ThirdPartyLink.Get() got length = %v, want %v", len(got), tt.wantLen)
				}
			},
		)
	}
}

func TestThirdPartyLink_List(t *testing.T) {
	common.RunListTest(t, `{
			"items": [
				{
					"linkingToken": "token123",
					"snippet": {"type": "channelToStoreLink"},
					"status": {"linkStatus": "linked"}
				}
			]
		}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			tpl := NewThirdPartyLink(
				WithService(svc),
				WithOutput(output),
				WithParts([]string{"snippet", "status"}),
			)
			return tpl.List
		},
	)
}

func TestThirdPartyLink_Insert(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "insert third party link",
			opts: []Option{
				WithLinkingToken("token123"),
				WithType("channelToStoreLink"),
				WithLinkStatus("pending"),
				WithParts([]string{"snippet", "status"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST, got %s", r.Method)
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					LinkingToken string `json:"linkingToken"`
					Snippet      struct {
						Type string `json:"type"`
					} `json:"snippet"`
					Status struct {
						LinkStatus string `json:"linkStatus"`
					} `json:"status"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.LinkingToken != "token123" {
					t.Errorf("expected linkingToken=token123, got %s", body.LinkingToken)
				}
				if body.Snippet.Type != "channelToStoreLink" {
					t.Errorf("expected snippet.type=channelToStoreLink, got %s", body.Snippet.Type)
				}
				if body.Status.LinkStatus != "pending" {
					t.Errorf("expected status.linkStatus=pending, got %s", body.Status.LinkStatus)
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
							_, _ = w.Write([]byte(`{"linkingToken": "token123"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				tpl := NewThirdPartyLink(opts...)
				var buf bytes.Buffer
				if err := tpl.Insert(&buf); (err != nil) != tt.wantErr {
					t.Errorf("ThirdPartyLink.Insert() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestThirdPartyLink_Update(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "update third party link",
			opts: []Option{
				WithLinkingToken("token123"),
				WithType("channelToStoreLink"),
				WithLinkStatus("linked"),
				WithParts([]string{"snippet", "status"}),
			},
			verify: func(r *http.Request) {
				if r.Method != "PUT" {
					return
				}

				defer func() { _ = r.Body.Close() }()
				var body struct {
					LinkingToken string `json:"linkingToken"`
					Status       struct {
						LinkStatus string `json:"linkStatus"`
					} `json:"status"`
				}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body.LinkingToken != "token123" {
					t.Errorf("expected linkingToken=token123, got %s", body.LinkingToken)
				}
				if body.Status.LinkStatus != "linked" {
					t.Errorf("expected status.linkStatus=linked, got %s", body.Status.LinkStatus)
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
							w.Header().Set("Content-Type", "application/json")
							if r.Method == "GET" {
								_, _ = w.Write([]byte(`{
									"items": [
										{"linkingToken": "token123", "snippet": {"type": "channelToStoreLink"}, "status": {"linkStatus": "pending"}}
									]
								}`))
								return
							}
							if tt.verify != nil {
								tt.verify(r)
							}
							_, _ = w.Write([]byte(`{"linkingToken": "token123"}`))
						},
					),
				)

				opts := append([]Option{WithService(svc)}, tt.opts...)
				tpl := NewThirdPartyLink(opts...)
				var buf bytes.Buffer
				if err := tpl.Update(&buf); (err != nil) != tt.wantErr {
					t.Errorf("ThirdPartyLink.Update() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestThirdPartyLink_Delete(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		verify  func(*http.Request)
		wantErr bool
	}{
		{
			name: "delete third party link",
			opts: []Option{
				WithLinkingToken("token123"),
				WithType("channelToStoreLink"),
			},
			verify: func(r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				if r.URL.Query().Get("linkingToken") != "token123" {
					t.Errorf("expected linkingToken=token123, got %s", r.URL.Query().Get("linkingToken"))
				}
				if r.URL.Query().Get("type") != "channelToStoreLink" {
					t.Errorf("expected type=channelToStoreLink, got %s", r.URL.Query().Get("type"))
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
				tpl := NewThirdPartyLink(opts...)
				var buf bytes.Buffer
				if err := tpl.Delete(&buf); (err != nil) != tt.wantErr {
					t.Errorf("ThirdPartyLink.Delete() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
