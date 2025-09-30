// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"encoding/base64"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

const (
	readTokenFailed  = "failed to read token"
	readSecretFailed = "failed to read client secret"
	authHint         = "Please configure client secret as described in https://github.com/eat-pray-ai/yutu#prerequisites"
)

type svc struct {
	Credential string `yaml:"credential" json:"credential"`
	CacheToken string `yaml:"cache_token" json:"cache_token"`
	credFile   string
	tokenFile  string

	service *youtube.Service
	ctx     context.Context
}

type Svc interface {
	GetService() *youtube.Service
	refreshClient() *http.Client
	newClient(*oauth2.Config) (*http.Client, *oauth2.Token)
	getConfig() *oauth2.Config
	startWebServer(string) chan string
	getTokenFromWeb(*oauth2.Config, string) *oauth2.Token
	getCodeFromPrompt(string) string
	saveToken(*oauth2.Token)
}

type Option func(*svc)

func NewY2BService(opts ...Option) Svc {
	s := &svc{}
	s.ctx = context.Background()
	s.credFile = "client_secret.json"

	for _, opt := range opts {
		opt(s)
	}
	return s
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
		// 2. cred is a base64 encoded string
		// 3. cred is a json string
		absCred, _ := filepath.Abs(cred)
		relCred, _ := filepath.Rel("/", absCred)

		if _, err := fs.Stat(fsys, relCred); err == nil {
			s.credFile = absCred
			credBytes, err := fs.ReadFile(fsys, relCred)
			if err != nil {
				slog.Error(
					readSecretFailed, "hint", authHint, "path", absCred, "error", err,
				)
				os.Exit(1)
			}
			s.Credential = string(credBytes)
			return
		}

		if credB64, err := base64.StdEncoding.DecodeString(cred); err == nil {
			s.Credential = string(credB64)
		} else if utils.IsJson(cred) {
			s.Credential = cred
		} else {
			slog.Error(parseSecretFailed, "hint", authHint, "error", err)
			os.Exit(1)
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
		// 2. token is a base64 encoded string
		// 3. token is a json string
		absToken, _ := filepath.Abs(token)
		relToken, _ := filepath.Rel("/", absToken)

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
			slog.Warn(parseTokenFailed, "error", err)
		}
	}
}
