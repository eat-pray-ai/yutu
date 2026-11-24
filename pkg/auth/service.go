// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/utils"
	"google.golang.org/api/youtube/v3"
)

const (
	authHint = "Please configure client secret as described in https://github.com/eat-pray-ai/yutu#prerequisites"

	readTokenFailed        = "failed to read token"
	readSecretFailed       = "failed to read client secret"
	defaultCallbackTimeout = 2 * time.Minute
)

type svc struct {
	Credential     string `yaml:"credential" json:"credential"`
	CacheToken     string `yaml:"cache_token" json:"cache_token"`
	credFile       string
	tokenFile      string
	redirectURL    string
	timeout        time.Duration
	initErr        error
	in             io.Reader
	out            io.Writer
	openURLFunc    func(string) error
	onAuthURL      func(string)
	codeProvider   func(ctx context.Context, authURL, redirectURL string) (string, error)
	noBrowser      bool
	nonInteractive bool

	service *youtube.Service
	ctx     context.Context
	state   string
}

type Svc interface {
	GetService() (*youtube.Service, error)
}

type Option func(*svc)

func NewY2BService(opts ...Option) Svc {
	s := &svc{}
	s.ctx = context.Background()
	s.credFile = "client_secret.json"
	s.timeout = defaultCallbackTimeout
	state, err := utils.RandomStage()
	if err != nil {
		s.initErr = err
	} else {
		s.state = state
	}
	s.in = os.Stdin
	s.out = os.Stdout

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithContext(ctx context.Context) Option {
	return func(s *svc) {
		s.ctx = ctx
	}
}

func WithRedirectURL(url string) Option {
	return func(s *svc) {
		s.redirectURL = url
	}
}

func WithCallbackTimeout(timeout time.Duration) Option {
	return func(s *svc) {
		s.timeout = timeout
	}
}

func WithIO(in io.Reader, out io.Writer) Option {
	return func(s *svc) {
		s.in = in
		if out == nil {
			out = io.Discard
		}
		s.out = out
	}
}

func WithNoBrowser(noBrowser bool) Option {
	return func(s *svc) {
		s.noBrowser = noBrowser
	}
}

func WithNonInteractive(nonInteractive bool) Option {
	return func(s *svc) {
		s.nonInteractive = nonInteractive
	}
}

func WithOpenURL(fn func(string) error) Option {
	return func(s *svc) {
		s.openURLFunc = fn
	}
}

func WithOnAuthURL(fn func(authURL string)) Option {
	return func(s *svc) {
		s.onAuthURL = fn
	}
}

func WithCodeProvider(
	fn func(ctx context.Context, authURL, redirectURL string) (string, error),
) Option {
	return func(s *svc) {
		s.codeProvider = fn
	}
}

func (s *svc) openURL(url string) (err error) {
	if s.openURLFunc != nil {
		return s.openURLFunc(url)
	}

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("cannot open URL %s on this platform", url)
	}

	return err
}

func WithCredential(cred string, fsys fs.FS) Option {
	return func(s *svc) {
		// cred > YUTU_CREDENTIAL
		envCred, ok := os.LookupEnv("YUTU_CREDENTIAL")
		if cred == "" && ok {
			cred = envCred
		} else if cred == "" {
			cred = s.credFile
		}
		// 1. cred is a file path
		// 2. cred is a Base64 encoded string
		// 3. cred is a JSON string
		absCred, _ := filepath.Abs(cred)
		relCred, _ := filepath.Rel(*pkg.RootDir, absCred)
		relCred = strings.ReplaceAll(relCred, `\`, `/`)

		if _, err := fs.Stat(fsys, relCred); err == nil {
			s.credFile = absCred
			credBytes, err := fs.ReadFile(fsys, relCred)
			if err != nil {
				s.initErr = fmt.Errorf("%s: %w", readSecretFailed, err)
				slog.Error(
					readSecretFailed, "hint", authHint, "path", absCred, "error", err,
				)
			}
			s.Credential = string(credBytes)
			return
		}

		if credB64, err := base64.StdEncoding.DecodeString(cred); err == nil {
			s.Credential = string(credB64)
		} else if utils.IsJson(cred) {
			s.Credential = cred
		} else {
			s.initErr = fmt.Errorf("%s: %w", parseSecretFailed, err)
			slog.Error(parseSecretFailed, "hint", authHint, "error", err)
		}
	}
}

func WithCacheToken(token string, fsys fs.FS) Option {
	return func(s *svc) {
		// token > YUTU_CACHE_TOKEN
		envToken, ok := os.LookupEnv("YUTU_CACHE_TOKEN")
		if token == "" && ok {
			token = envToken
		} else if token == "" {
			token = "youtube.token.json"
		}

		// 1. token is a file path
		// 2. token is a Base64 encoded string
		// 3. token is a JSON string
		absToken, _ := filepath.Abs(token)
		relToken, _ := filepath.Rel(*pkg.RootDir, absToken)
		relToken = strings.ReplaceAll(relToken, `\`, `/`)

		if _, err := fs.Stat(fsys, relToken); err == nil {
			tokenBytes, err := fs.ReadFile(fsys, relToken)
			if err != nil {
				slog.Warn(readTokenFailed, "path", absToken, "error", err)
			}
			s.tokenFile = relToken
			s.CacheToken = string(tokenBytes)
			return
		} else if os.IsNotExist(err) && strings.HasSuffix(token, ".json") {
			s.tokenFile = relToken
			return
		}

		if tokenB64, err := base64.StdEncoding.DecodeString(token); err == nil {
			s.CacheToken = string(tokenB64)
		} else if utils.IsJson(token) {
			s.CacheToken = token
		} else {
			slog.Warn(parseTokenFailed, "token", token, "error", err)
		}
	}
}
