// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// NewTestService creates a youtube.Service backed by the given handler for testing.
// It registers cleanup of the test server automatically.
func NewTestService(t *testing.T, handler http.Handler) *youtube.Service {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	svc, err := youtube.NewService(
		context.Background(),
		option.WithEndpoint(ts.URL),
		option.WithAPIKey("test-key"),
	)
	if err != nil {
		t.Fatalf("failed to create youtube service: %v", err)
	}
	return svc
}

// PaginationHandler returns an http.HandlerFunc that simulates 2-page pagination.
// idPrefix is used for item IDs (e.g., "channel", "video").
// An optional itemFmt function can override the default item JSON template;
// it receives the idPrefix and index, and returns the JSON for one item.
func PaginationHandler(idPrefix string, itemFmt ...func(string, int) string) http.HandlerFunc {
	fmtItem := func(prefix string, i int) string {
		return fmt.Sprintf(`{"id": "%s-%d"}`, prefix, i)
	}
	if len(itemFmt) > 0 && itemFmt[0] != nil {
		fmtItem = itemFmt[0]
	}
	return func(w http.ResponseWriter, r *http.Request) {
		pageToken := r.URL.Query().Get("pageToken")
		w.Header().Set("Content-Type", "application/json")
		if pageToken == "" {
			items := make([]string, 20)
			for i := range 20 {
				items[i] = fmtItem(idPrefix, i)
			}
			_, _ = w.Write(
				fmt.Appendf(
					nil,
					`{"items": [%s], "nextPageToken": "page-2"}`,
					strings.Join(items, ","),
				),
			)
		} else if pageToken == "page-2" {
			_, _ = w.Write(
				[]byte(
					fmt.Sprintf(
						`{"items": [%s, %s], "nextPageToken": ""}`,
						fmtItem(idPrefix, 20), fmtItem(idPrefix, 21),
					),
				),
			)
		}
	}
}

// RunListTest runs the standard json/yaml/table output test matrix.
// mockResponse is the JSON the mock server returns.
// listFn receives a *youtube.Service and output format, returns a function
// that calls List on the resource with a writer.
func RunListTest(
	t *testing.T,
	mockResponse string,
	listFn func(svc *youtube.Service, output string) func(io.Writer) error,
) {
	t.Helper()
	svc := NewTestService(
		t, http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(mockResponse))
			},
		),
	)

	for _, output := range []string{"json", "yaml", "table"} {
		t.Run(
			"list "+output, func(t *testing.T) {
				var buf bytes.Buffer
				if err := listFn(svc, output)(&buf); err != nil {
					t.Errorf("List(%s) error = %v", output, err)
				}
				if buf.Len() == 0 {
					t.Errorf("List(%s) output is empty", output)
				}
			},
		)
	}
}
