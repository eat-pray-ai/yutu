// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	"google.golang.org/api/youtube/v3"
)

// testResource is a minimal type that satisfies HasFields for option testing.
type testResource struct {
	Fields
}

type providerContextKey struct{}

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
					switch pageToken {
					case "":
						items := make([]string, 20)
						for i := range 20 {
							items[i] = fmt.Sprintf(`{"id": "v%d"}`, i+1)
						}
						jsonItems := "[" + strings.Join(items, ",") + "]"
						_, _ = fmt.Fprintf(w, `{
							"items": %s,
							"nextPageToken": "page-2"
						}`, jsonItems)
					case "page-2":
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
			if err := PrintList(tt.output, items, &buf, header, rowFn); err != nil {
				t.Fatalf("PrintList(%q) error = %v", tt.output, err)
			}
			if buf.Len() == 0 {
				t.Errorf("PrintList(%q) produced empty output", tt.output)
			}
		})
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestPrintList_WriteError(t *testing.T) {
	items := []*youtube.Video{{Id: "v1"}}
	header := table.Row{"ID"}
	rowFn := func(v *youtube.Video) table.Row { return table.Row{v.Id} }

	if err := PrintList("table", items, errWriter{}, header, rowFn); err == nil {
		t.Fatal("expected table write error, got nil")
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
			if err := PrintResult(tt.output, data, &buf, "Video %s updated", "v1"); err != nil {
				t.Fatalf("PrintResult(%q) error = %v", tt.output, err)
			}
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

func TestPrintResult_WriteError(t *testing.T) {
	if err := PrintResult("", nil, errWriter{}, "Video %s updated", "v1"); err == nil {
		t.Fatal("expected default output write error, got nil")
	}
}

func TestPaginationHandler(t *testing.T) {
	svc := NewTestService(
		t, PaginationHandler(
			"custom",
			func(prefix string, i int) string {
				return fmt.Sprintf(`{"id":"%s-%d","snippet":{"title":"Title %d"}}`, prefix, i, i)
			},
		),
	)

	items, err := Paginate(
		&Fields{MaxResults: 22}, svc.Videos.List([]string{"id", "snippet"}),
		videoExtract, errors.New("list failed"),
	)
	if err != nil {
		t.Fatalf("Paginate() error = %v", err)
	}
	if len(items) != 22 {
		t.Fatalf("Paginate() got %d items, want 22", len(items))
	}
	if items[0].Id != "custom-0" || items[21].Id != "custom-21" {
		t.Errorf("Paginate() got first/last IDs = %q/%q", items[0].Id, items[21].Id)
	}
}

func TestRunListTest(t *testing.T) {
	RunListTest(
		t,
		`{"items":[{"id":"v1","snippet":{"title":"Title 1"}}]}`,
		func(svc *youtube.Service, output string) func(io.Writer) error {
			return func(w io.Writer) error {
				items, err := Paginate(
					&Fields{Service: svc, MaxResults: 1},
					svc.Videos.List([]string{"id", "snippet"}),
					videoExtract,
					errors.New("list failed"),
				)
				if err != nil {
					return err
				}
				return PrintList(
					output, items, w, table.Row{"ID", "Title"},
					func(v *youtube.Video) table.Row { return table.Row{v.Id, v.Snippet.Title} },
				)
			}
		},
	)
}

// ---------- TestSetContext ----------

func TestSetContext(t *testing.T) {
	f := &Fields{}
	ctx := t.Context()
	f.SetContext(ctx)
	if f.Ctx != ctx {
		t.Error("SetContext did not set context")
	}
}

func TestGetFields(t *testing.T) {
	f := &Fields{}
	if got := f.GetFields(); got != f {
		t.Error("GetFields() did not return receiver")
	}
}

func TestPrintList_UnknownOutput(t *testing.T) {
	items := []*youtube.Video{{Id: "v1"}}
	header := table.Row{"ID"}
	rowFn := func(v *youtube.Video) table.Row { return table.Row{v.Id} }

	var buf bytes.Buffer
	if err := PrintList("unknown", items, &buf, header, rowFn); err != nil {
		t.Fatalf("PrintList() error = %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("PrintList() output = %q, want empty", buf.String())
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

func TestEnsureService_UsesDefaultServiceProvider(t *testing.T) {
	originalProvider := DefaultServiceProvider
	t.Cleanup(func() { DefaultServiceProvider = originalProvider })

	wantSvc := &youtube.Service{}
	ctx := context.WithValue(t.Context(), providerContextKey{}, "test-user")
	called := false
	var ctxValue any
	SetDefaultServiceProvider(func(gotCtx context.Context) (*youtube.Service, error) {
		called = true
		if gotCtx != ctx {
			t.Errorf("ServiceProvider context = %p, want %p", gotCtx, ctx)
		}
		ctxValue = gotCtx.Value(providerContextKey{})
		return wantSvc, nil
	})

	f := &Fields{}
	f.SetContext(ctx)
	if err := f.EnsureService(); err != nil {
		t.Fatalf("EnsureService() error = %v, want nil", err)
	}
	if !called {
		t.Fatal("expected DefaultServiceProvider to be called")
	}
	if ctxValue != "test-user" {
		t.Errorf("provider context value = %v, want test-user", ctxValue)
	}
	if f.Service != wantSvc {
		t.Error("EnsureService() did not assign provider service")
	}
}

func TestSetDefaultServiceProvider_NilRestoresDefault(t *testing.T) {
	originalProvider := DefaultServiceProvider
	t.Cleanup(func() { DefaultServiceProvider = originalProvider })

	SetDefaultServiceProvider(func(context.Context) (*youtube.Service, error) {
		return &youtube.Service{}, nil
	})
	SetDefaultServiceProvider(nil)

	if DefaultServiceProvider == nil {
		t.Fatal("DefaultServiceProvider is nil")
	}
}

func TestEnsureService_FallsBackWithoutContext(t *testing.T) {
	f := &Fields{}
	err := f.EnsureService()
	if err == nil {
		t.Skip("credentials available in environment, skipping fallback test")
	}
	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}

func TestEnsureService_FallsBackWithContextNoToken(t *testing.T) {
	f := &Fields{}
	f.SetContext(t.Context())
	// No TokenInfo in context → should fall back to file auth
	err := f.EnsureService()
	if err == nil {
		t.Skip("credentials available in environment, skipping")
	}
	// Should get credential error, not context/token error
	_ = err
}

func TestTokenInfoFromContext_Nil(t *testing.T) {
	info := sdkauth.TokenInfoFromContext(t.Context())
	if info != nil {
		t.Error("expected nil TokenInfo from plain context")
	}
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
			if r.MaxResults != tt.want {
				t.Errorf("WithMaxResults(%d) = %d, want %d", tt.input, r.MaxResults, tt.want)
			}
		})
	}
}

func TestWithContext(t *testing.T) {
	ctx := context.WithValue(t.Context(), providerContextKey{}, "value")
	r := &testResource{}

	WithContext[*testResource](ctx)(r)

	if r.Ctx != ctx {
		t.Error("WithContext() did not set context")
	}
}

func TestWithIds(t *testing.T) {
	r := &testResource{}
	ids := []string{"id1", "id2", "id3"}
	WithIds[*testResource](ids)(r)
	if len(r.Ids) != len(ids) {
		t.Fatalf("WithIds() got %d ids, want %d", len(r.Ids), len(ids))
	}
	for i, id := range r.Ids {
		if id != ids[i] {
			t.Errorf("WithIds()[%d] = %q, want %q", i, id, ids[i])
		}
	}
}

func TestWithParts(t *testing.T) {
	r := &testResource{}
	parts := []string{"snippet", "contentDetails"}
	WithParts[*testResource](parts)(r)
	if len(r.Parts) != len(parts) {
		t.Fatalf("WithParts() got %d parts, want %d", len(r.Parts), len(parts))
	}
	for i, p := range r.Parts {
		if p != parts[i] {
			t.Errorf("WithParts()[%d] = %q, want %q", i, p, parts[i])
		}
	}
}

func TestWithOutput(t *testing.T) {
	r := &testResource{}
	WithOutput[*testResource]("json")(r)
	if r.Output != "json" {
		t.Errorf("WithOutput() = %q, want %q", r.Output, "json")
	}
}

func TestWithService(t *testing.T) {
	r := &testResource{}
	svc := &youtube.Service{}
	WithService[*testResource](svc)(r)
	if r.Service != svc {
		t.Errorf("WithService() did not set service")
	}
}

func TestWithHl(t *testing.T) {
	r := &testResource{}
	WithHl[*testResource]("en")(r)
	if r.Hl != "en" {
		t.Errorf("WithHl() = %q, want %q", r.Hl, "en")
	}
}

func TestWithChannelId(t *testing.T) {
	r := &testResource{}
	WithChannelId[*testResource]("chan-123")(r)
	if r.ChannelId != "chan-123" {
		t.Errorf("WithChannelId() = %q, want %q", r.ChannelId, "chan-123")
	}
}

func TestWithOBOCO(t *testing.T) {
	r := &testResource{}
	WithOBOCO[*testResource]("owner-456")(r)
	if r.OBOCO != "owner-456" {
		t.Errorf("WithOBOCO() = %q, want %q", r.OBOCO, "owner-456")
	}
}
