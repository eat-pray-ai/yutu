// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	_ "embed"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eat-pray-ai/yutu/pkg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	createSvcFailed    = "failed to create YouTube service"
	parseUrlFailed     = "failed to parse redirect URL"
	stateMatchFailed   = "state doesn't match, possible CSRF attack"
	readPromptFailed   = "failed to read prompt"
	exchangeFailed     = "failed to exchange token"
	listenFailed       = "failed to start web server"
	cacheTokenFailed   = "failed to cache token"
	parseTokenFailed   = "failed to parse token"
	refreshTokenFailed = "failed to refresh token, please re-authenticate in cli"
	parseSecretFailed  = "failed to parse client secret"
	authTimeout        = "authorization timed out after %s"

	browserOpenedHint = "Your browser has been opened to an authorization URL. yutu will resume once authorization has been provided.\n%s\n"
	openBrowserHint   = "It seems that your browser is not open. Go to the following link in your browser:\n%s\n"
	manualInputHint   = `
After completing the authorization flow, enter the authorization code on command line.

If you end up in an error page after completing the authorization flow,
and the url in the address bar is in the form of
'%s/?state=DONOT-COPY&code=COPY-THIS&scope=DONOT-COPY'
ONLY 'COPY-THIS' part is the code you need to enter on command line.
`
)

// authorizedPage is shown in the browser once the OAuth redirect is received.
// It must stay self-contained: the loopback server closes right after serving it.
//
//go:embed authorized.html
var authorizedPage []byte

// Scopes is the set of YouTube OAuth scopes required by yutu.
var Scopes = []string{
	youtube.YoutubeScope,
	youtube.YoutubeForceSslScope,
	youtube.YoutubeChannelMembershipsCreatorScope,
}

func (s *svc) GetService() (*youtube.Service, error) {
	if s.initErr != nil {
		return nil, s.initErr
	}

	client, err := s.refreshClient()
	if err != nil {
		return nil, err
	}
	service, err := youtube.NewService(s.ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", createSvcFailed, err)
	}
	s.service = service

	return s.service, nil
}

func (s *svc) refreshClient() (client *http.Client, err error) {
	config, err := s.getConfig()
	if err != nil {
		return nil, err
	}
	authedToken := &oauth2.Token{}
	err = json.Unmarshal([]byte(s.CacheToken), authedToken)
	if err != nil {
		client, authedToken, err = s.newClient(config)
		if err != nil {
			return nil, err
		}
		if s.tokenFile != "" {
			if err := s.saveToken(authedToken); err != nil {
				return nil, err
			}
		}
		return client, nil
	}

	if !authedToken.Valid() {
		tokenSource := config.TokenSource(s.ctx, authedToken)
		authedToken, err = tokenSource.Token()
		if err != nil && s.tokenFile != "" {
			client, authedToken, err = s.newClient(config)
			if err != nil {
				return nil, err
			}
			if err := s.saveToken(authedToken); err != nil {
				return nil, err
			}
			return client, nil
		} else if err != nil {
			return nil, fmt.Errorf("%s: %w", refreshTokenFailed, err)
		}

		if authedToken != nil && s.tokenFile != "" {
			if err := s.saveToken(authedToken); err != nil {
				return nil, err
			}
		}
		return config.Client(s.ctx, authedToken), nil
	}

	return config.Client(s.ctx, authedToken), nil
}

func (s *svc) newClient(config *oauth2.Config) (
	client *http.Client, token *oauth2.Token, err error,
) {
	verifier := oauth2.GenerateVerifier()
	authURL := config.AuthCodeURL(
		s.state,
		oauth2.ApprovalForce,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)
	token, err = s.getTokenFromWeb(config, authURL, verifier)
	if err != nil {
		return nil, nil, err
	}
	client = config.Client(s.ctx, token)

	return client, token, nil
}

func (s *svc) getConfig() (*oauth2.Config, error) {
	config, err := google.ConfigFromJSON([]byte(s.Credential), Scopes...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parseSecretFailed, err)
	}

	if s.redirectURL != "" {
		config.RedirectURL = s.redirectURL
	}

	return config, nil
}

func (s *svc) startWebServer(redirectURL string) (chan string, func() error, error) {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", parseUrlFailed, err)
	}

	listener, err := net.Listen("tcp", u.Host)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", listenFailed, err)
	}

	codeCh := make(chan string, 1)
	go func() {
		_ = http.Serve(
			listener, http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/" {
						http.NotFound(w, r)
						return
					}
					state := r.FormValue("state")
					if state != s.state {
						slog.Error(stateMatchFailed)
						http.Error(w, stateMatchFailed, http.StatusBadRequest)
						return
					}
					code := r.FormValue("code")
					select {
					case codeCh <- code:
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
						_, _ = w.Write(authorizedPage)
						_ = listener.Close()
					default:
						http.Error(
							w, "authorization code already received",
							http.StatusConflict,
						)
					}
				},
			),
		)
	}()

	return codeCh, listener.Close, nil
}

func (s *svc) waitForCode(codeCh <-chan string) (string, error) {
	select {
	case code := <-codeCh:
		return code, nil
	case <-s.ctx.Done():
		return "", s.ctx.Err()
	case <-time.After(s.timeout):
		return "", fmt.Errorf(authTimeout, s.timeout)
	}
}

func (s *svc) getCodeFromPrompt(authURL string, redirectURL string) (code string, err error) {
	_, _ = fmt.Fprintf(s.out, openBrowserHint, authURL)
	_, _ = fmt.Fprintf(s.out, manualInputHint, redirectURL)
	_, err = fmt.Fscan(s.in, &code)
	if err != nil {
		return "", fmt.Errorf("%s: %w", readPromptFailed, err)
	}

	if strings.HasPrefix(code, "4%2F") {
		code = strings.Replace(code, "4%2F", "4/", 1)
	}
	return code, nil
}

func (s *svc) getTokenFromWeb(
	config *oauth2.Config, authURL string, verifier string,
) (*oauth2.Token, error) {
	codeCh, closeServer, err := s.startWebServer(config.RedirectURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = closeServer() }()

	var code string
	if err := s.openURL(authURL); err == nil {
		_, _ = fmt.Fprintf(s.out, browserOpenedHint, authURL)
		code, err = s.waitForCode(codeCh)
		if err != nil {
			return nil, err
		}
	}

	if code == "" {
		code, err = s.getCodeFromPrompt(authURL, config.RedirectURL)
		if err != nil {
			return nil, err
		}
	}

	slog.Debug("Authorization code received")
	token, err := config.Exchange(
		s.ctx, code, oauth2.VerifierOption(verifier),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", exchangeFailed, err)
	}

	return token, nil
}

func (s *svc) saveToken(token *oauth2.Token) error {
	dir := filepath.Dir(s.tokenFile)
	if err := pkg.Root.MkdirAll(dir, 0755); err != nil {
		slog.Error(cacheTokenFailed, "dir", dir, "error", err)
		return fmt.Errorf("%s: %w", cacheTokenFailed, err)
	}

	f, err := pkg.Root.OpenFile(
		s.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600,
	)
	if err != nil {
		slog.Error(cacheTokenFailed, "file", s.tokenFile, "error", err)
		return fmt.Errorf("%s: %w", cacheTokenFailed, err)
	}
	defer func() { _ = f.Close() }()
	err = json.MarshalWrite(f, token)
	if err != nil {
		return fmt.Errorf("%s: %w", cacheTokenFailed, err)
	}
	slog.Debug("Token cached to file", "file", s.tokenFile)

	return nil
}
