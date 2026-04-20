package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
)

func TestGoogleTokenVerifier_ValidToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("access_token")
		if token != "valid-google-token" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_token"})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"expires_in": "3600",
			"scope":      "https://www.googleapis.com/auth/youtube https://www.googleapis.com/auth/youtube.force-ssl https://www.googleapis.com/auth/youtube.channel-memberships.creator",
			"sub":        "user-123",
		})
	}))
	defer ts.Close()

	verifier := NewGoogleTokenVerifier(ts.URL)
	info, err := verifier(context.Background(), "valid-google-token", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.UserID != "user-123" {
		t.Errorf("UserID = %q, want 'user-123'", info.UserID)
	}
	if len(info.Scopes) != 3 {
		t.Errorf("Scopes = %v, want 3 scopes", info.Scopes)
	}
	raw, ok := info.Extra["access_token"]
	if !ok || raw != "valid-google-token" {
		t.Errorf("Extra[access_token] = %v, want 'valid-google-token'", raw)
	}
	if info.Expiration.IsZero() {
		t.Error("expected non-zero expiration")
	}
}

func TestGoogleTokenVerifier_InvalidToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid_token"})
	}))
	defer ts.Close()

	verifier := NewGoogleTokenVerifier(ts.URL)
	_, err := verifier(context.Background(), "bad-token", nil)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if !errors.Is(err, sdkauth.ErrInvalidToken) {
		t.Errorf("error = %v, want ErrInvalidToken", err)
	}
}

func TestGoogleTokenVerifier_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	verifier := NewGoogleTokenVerifier(ts.URL)
	_, err := verifier(context.Background(), "any-token", nil)
	if err == nil {
		t.Fatal("expected error for server failure")
	}
	if !errors.Is(err, sdkauth.ErrInvalidToken) {
		t.Errorf("error = %v, want ErrInvalidToken", err)
	}
}
