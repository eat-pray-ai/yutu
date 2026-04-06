// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"net/http"
	"net/http/httptest"
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
