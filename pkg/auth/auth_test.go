package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eat-pray-ai/yutu/pkg"
	"golang.org/x/oauth2"
)

// helper to build a minimal valid Google OAuth2 credential JSON.
func validCredentialJSON(redirectURL string) string {
	cred := map[string]map[string]any{
		"web": {
			"client_id":                   "test-client-id",
			"project_id":                  "test-project",
			"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
			"token_uri":                   "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret":               "test-secret",
			"redirect_uris":               []string{redirectURL},
		},
	}
	b, _ := json.Marshal(cred)
	return string(b)
}

func TestGetConfig_Success(t *testing.T) {
	s := NewY2BService(
		WithCredential(validCredentialJSON("http://localhost:8216"), os.DirFS(".")),
	).(*svc)

	config, err := s.getConfig()
	if err != nil {
		t.Fatalf("getConfig returned error: %v", err)
	}

	if config.ClientID != "test-client-id" {
		t.Errorf(
			"unexpected ClientID: got %q, want %q", config.ClientID, "test-client-id",
		)
	}
	if config.ClientSecret != "test-secret" {
		t.Errorf(
			"unexpected ClientSecret: got %q, want %q", config.ClientSecret,
			"test-secret",
		)
	}
	if config.RedirectURL != "http://localhost:8216" {
		t.Errorf(
			"unexpected RedirectURL: got %q, want %q", config.RedirectURL,
			"http://localhost:8216",
		)
	}
	if len(config.Scopes) == 0 {
		t.Fatalf("expected scopes to be populated, got none")
	}
}

func TestGetConfig_InvalidJSON(t *testing.T) {
	s := NewY2BService(
		WithCredential("not-a-json", os.DirFS(".")),
	).(*svc)

	_, err := s.getConfig()
	if err == nil {
		t.Fatalf("expected error from getConfig with invalid JSON, got nil")
	}
	if err.Error() == "" {
		t.Fatalf("expected non-empty error message")
	}
}

func TestStartWebServer_InvalidURL(t *testing.T) {
	s := NewY2BService().(*svc)
	s.state = "state"

	_, err := s.startWebServer("://bad-url")
	if err == nil {
		t.Fatalf("expected error for invalid URL, got nil")
	}
	if err.Error() == "" {
		t.Fatalf("expected non-empty error message")
	}
}

func TestStartWebServer_PortConflict(t *testing.T) {
	// Occupy a TCP port first.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on random port: %v", err)
	}
	defer func(ln net.Listener) {
		_ = ln.Close()
	}(ln)

	addr := ln.Addr().String()
	redirectURL := fmt.Sprintf("http://%s", addr)

	s := NewY2BService().(*svc)
	s.state = "state"

	_, err = s.startWebServer(redirectURL)
	if err == nil {
		t.Fatalf("expected error due to port conflict, got nil")
	}
}

func TestStartWebServer_StateMismatch(t *testing.T) {
	// Use an ephemeral port so startWebServer succeeds.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on random port: %v", err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()

	redirectURL := fmt.Sprintf("http://%s", addr)

	s := NewY2BService().(*svc)
	s.state = "expected-state"
	codeCh, err := s.startWebServer(redirectURL)
	if err != nil {
		t.Fatalf("startWebServer returned error: %v", err)
	}

	// Send a request with mismatched state.
	resp, err := http.Get(redirectURL + "/?state=wrong&code=code123")
	if err != nil {
		t.Fatalf("failed to send HTTP request: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case <-codeCh:
		t.Fatalf("expected no code to be sent on state mismatch")
	case <-time.After(100 * time.Millisecond):
		// ok: channel should remain empty
	}
}

func TestGetCodeFromPrompt_Success(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("some-code")

	s := NewY2BService(WithIO(in, &out)).(*svc)

	code, err := s.getCodeFromPrompt("http://example.com/auth", "http://localhost:8216")
	if err != nil {
		t.Fatalf("getCodeFromPrompt returned error: %v", err)
	}
	if code != "some-code" {
		t.Fatalf("unexpected code: got %q, want %q", code, "some-code")
	}

	outStr := out.String()
	if !strings.Contains(outStr, openBrowserHint[:20]) { // basic sanity check that hint was written
		t.Fatalf("expected openBrowserHint to be written to out, got %q", outStr)
	}
	if !strings.Contains(outStr, "After completing the authorization flow") {
		t.Fatalf("expected manualInputHint to be written to out, got %q", outStr)
	}
	if !strings.Contains(outStr, "http://localhost:8216/?state=DONOT-COPY") {
		t.Fatalf("expected redirect URL in manualInputHint, got %q", outStr)
	}
}

func TestGetCodeFromPrompt_URLDecode(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("4%2Fabc")

	s := NewY2BService(WithIO(in, &out)).(*svc)

	code, err := s.getCodeFromPrompt("http://example.com/auth", "http://localhost:8216")
	if err != nil {
		t.Fatalf("getCodeFromPrompt returned error: %v", err)
	}
	if code != "4/abc" {
		t.Fatalf("unexpected code: got %q, want %q", code, "4/abc")
	}
}

func TestSaveToken_InvalidPath(t *testing.T) {
	// Point pkg.Root to a directory without write permission so that
	// the underlying filesystem operations fail.
	tmpDir := t.TempDir()
	if err := os.Chmod(tmpDir, 0o555); err != nil {
		t.Fatalf("failed to chmod temp dir: %v", err)
	}

	origRootDir := *pkg.RootDir
	origRoot := pkg.Root
	defer func() {
		pkg.RootDir = &origRootDir
		pkg.Root = origRoot
	}()

	pkg.RootDir = &tmpDir
	root, err := os.OpenRoot(tmpDir)
	if err != nil {
		t.Fatalf("failed to open root: %v", err)
	}
	pkg.Root = root

	s := NewY2BService().(*svc)
	s.tokenFile = "token.json"
	tok := &oauth2.Token{AccessToken: "abc"}

	err = s.saveToken(tok)
	if err == nil {
		t.Fatalf("expected error from saveToken with invalid path, got nil")
	}
}
