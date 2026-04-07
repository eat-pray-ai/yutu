// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/api/youtube/v3"
)

// testResource is a minimal type that satisfies HasFields for option testing.
type testResource struct {
	Fields
}

func (r *testResource) GetFields() *Fields {
	return &r.Fields
}

func (r *testResource) EnsureService() error {
	return r.Fields.EnsureService()
}

// videoExtract is the extract function for Paginate with youtube.VideoListResponse.
func videoExtract(r *youtube.VideoListResponse) ([]*youtube.Video, string) {
	return r.Items, r.NextPageToken
}

// ---------- TestPaginate ----------

func TestPaginate(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		maxResults int64
		wantLen    int
		wantErr    bool
	}{
		{
			name: "single page",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"items": [
						{"id": "v1"},
						{"id": "v2"}
					],
					"nextPageToken": ""
				}`))
			}),
			maxResults: 10,
			wantLen:    2,
			wantErr:    false,
		},
		{
			name: "multi page",
			handler: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					pageToken := r.URL.Query().Get("pageToken")
					w.Header().Set("Content-Type", "application/json")
					if pageToken == "" {
						items := make([]string, 20)
						for i := range 20 {
							items[i] = fmt.Sprintf(`{"id": "v%d"}`, i+1)
						}
						jsonItems := "[" + strings.Join(items, ",") + "]"
						_, _ = fmt.Fprintf(w, `{
							"items": %s,
							"nextPageToken": "page-2"
						}`, jsonItems)
					} else if pageToken == "page-2" {
						_, _ = w.Write([]byte(`{
							"items": [{"id": "v21"}, {"id": "v22"}],
							"nextPageToken": ""
						}`))
					}
				}
			}(),
			maxResults: 22,
			wantLen:    22,
			wantErr:    false,
		},
		{
			name: "zero results",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"items": [],
					"nextPageToken": ""
				}`))
			}),
			maxResults: 10,
			wantLen:    0,
			wantErr:    false,
		},
		{
			name: "API error first page",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}),
			maxResults: 10,
			wantLen:    0,
			wantErr:    true,
		},
		{
			name: "API error second page",
			handler: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					pageToken := r.URL.Query().Get("pageToken")
					w.Header().Set("Content-Type", "application/json")
					if pageToken == "" {
						items := make([]string, 20)
						for i := range 20 {
							items[i] = fmt.Sprintf(`{"id": "v%d"}`, i+1)
						}
						jsonItems := "[" + strings.Join(items, ",") + "]"
						_, _ = fmt.Fprintf(w, `{
							"items": %s,
							"nextPageToken": "page-2"
						}`, jsonItems)
					} else {
						http.Error(w, "internal server error", http.StatusInternalServerError)
					}
				}
			}(),
			maxResults: 40,
			wantLen:    20,
			wantErr:    true,
		},
		{
			name: "maxResults smaller than page",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				items := make([]string, 5)
				for i := range 5 {
					items[i] = fmt.Sprintf(`{"id": "v%d"}`, i+1)
				}
				jsonItems := "[" + strings.Join(items, ",") + "]"
				_, _ = fmt.Fprintf(w, `{
					"items": %s,
					"nextPageToken": ""
				}`, jsonItems)
			}),
			maxResults: 5,
			wantLen:    5,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTestService(t, tt.handler)
			f := &Fields{
				Service:    svc,
				MaxResults: tt.maxResults,
			}
			call := svc.Videos.List([]string{"id"})
			got, err := Paginate(f, call, videoExtract, fmt.Errorf("video list failed"))
			if (err != nil) != tt.wantErr {
				t.Errorf("Paginate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != tt.wantLen {
				t.Errorf("Paginate() got %d items, want %d", len(got), tt.wantLen)
			}
		})
	}
}

// ---------- TestPrintList ----------

func TestPrintList(t *testing.T) {
	items := []*youtube.Video{
		{Id: "v1", Snippet: &youtube.VideoSnippet{Title: "Title 1"}},
		{Id: "v2", Snippet: &youtube.VideoSnippet{Title: "Title 2"}},
	}
	header := table.Row{"ID", "Title"}
	rowFn := func(v *youtube.Video) table.Row {
		return table.Row{v.Id, v.Snippet.Title}
	}

	tests := []struct {
		name   string
		output string
	}{
		{name: "json output", output: "json"},
		{name: "yaml output", output: "yaml"},
		{name: "table output", output: "table"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintList(tt.output, items, &buf, header, rowFn)
			if buf.Len() == 0 {
				t.Errorf("PrintList(%q) produced empty output", tt.output)
			}
		})
	}
}

// ---------- TestPrintResult ----------

func TestPrintResult(t *testing.T) {
	data := map[string]string{"id": "v1", "title": "Title 1"}

	tests := []struct {
		name      string
		output    string
		wantEmpty bool
		wantSub   string
	}{
		{
			name:      "json output",
			output:    "json",
			wantEmpty: false,
		},
		{
			name:      "yaml output",
			output:    "yaml",
			wantEmpty: false,
		},
		{
			name:      "silent output",
			output:    "silent",
			wantEmpty: true,
		},
		{
			name:      "default output",
			output:    "",
			wantEmpty: false,
			wantSub:   "Video v1 updated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintResult(tt.output, data, &buf, "Video %s updated", "v1")
			if tt.wantEmpty && buf.Len() != 0 {
				t.Errorf("PrintResult(%q) should produce no output, got %q", tt.output, buf.String())
			}
			if !tt.wantEmpty && buf.Len() == 0 {
				t.Errorf("PrintResult(%q) produced empty output", tt.output)
			}
			if tt.wantSub != "" && !strings.Contains(buf.String(), tt.wantSub) {
				t.Errorf("PrintResult(%q) = %q, want substring %q", tt.output, buf.String(), tt.wantSub)
			}
		})
	}
}

// ---------- TestEnsureService ----------

func TestEnsureService(t *testing.T) {
	t.Run("service already set", func(t *testing.T) {
		svc := &youtube.Service{}
		f := &Fields{Service: svc}
		err := f.EnsureService()
		if err != nil {
			t.Errorf("EnsureService() error = %v, want nil", err)
		}
		if f.Service != svc {
			t.Errorf("EnsureService() changed Service pointer")
		}
	})
}

// ---------- Generic option functions ----------

func TestWithMaxResults(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  int64
	}{
		{name: "positive", input: 42, want: 42},
		{name: "zero becomes MaxInt64", input: 0, want: math.MaxInt64},
		{name: "negative becomes 1", input: -5, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &testResource{}
			WithMaxResults[*testResource](tt.input)(r)
			if r.Fields.MaxResults != tt.want {
				t.Errorf("WithMaxResults(%d) = %d, want %d", tt.input, r.Fields.MaxResults, tt.want)
			}
		})
	}
}

func TestWithIds(t *testing.T) {
	r := &testResource{}
	ids := []string{"id1", "id2", "id3"}
	WithIds[*testResource](ids)(r)
	if len(r.Fields.Ids) != len(ids) {
		t.Fatalf("WithIds() got %d ids, want %d", len(r.Fields.Ids), len(ids))
	}
	for i, id := range r.Fields.Ids {
		if id != ids[i] {
			t.Errorf("WithIds()[%d] = %q, want %q", i, id, ids[i])
		}
	}
}

func TestWithParts(t *testing.T) {
	r := &testResource{}
	parts := []string{"snippet", "contentDetails"}
	WithParts[*testResource](parts)(r)
	if len(r.Fields.Parts) != len(parts) {
		t.Fatalf("WithParts() got %d parts, want %d", len(r.Fields.Parts), len(parts))
	}
	for i, p := range r.Fields.Parts {
		if p != parts[i] {
			t.Errorf("WithParts()[%d] = %q, want %q", i, p, parts[i])
		}
	}
}

func TestWithOutput(t *testing.T) {
	r := &testResource{}
	WithOutput[*testResource]("json")(r)
	if r.Fields.Output != "json" {
		t.Errorf("WithOutput() = %q, want %q", r.Fields.Output, "json")
	}
}

func TestWithService(t *testing.T) {
	r := &testResource{}
	svc := &youtube.Service{}
	WithService[*testResource](svc)(r)
	if r.Fields.Service != svc {
		t.Errorf("WithService() did not set service")
	}
}

func TestWithHl(t *testing.T) {
	r := &testResource{}
	WithHl[*testResource]("en")(r)
	if r.Fields.Hl != "en" {
		t.Errorf("WithHl() = %q, want %q", r.Fields.Hl, "en")
	}
}

func TestWithChannelId(t *testing.T) {
	r := &testResource{}
	WithChannelId[*testResource]("chan-123")(r)
	if r.Fields.ChannelId != "chan-123" {
		t.Errorf("WithChannelId() = %q, want %q", r.Fields.ChannelId, "chan-123")
	}
}

func TestWithOnBehalfOfContentOwner(t *testing.T) {
	r := &testResource{}
	WithOnBehalfOfContentOwner[*testResource]("owner-456")(r)
	if r.Fields.OnBehalfOfContentOwner != "owner-456" {
		t.Errorf("WithOnBehalfOfContentOwner() = %q, want %q", r.Fields.OnBehalfOfContentOwner, "owner-456")
	}
}
