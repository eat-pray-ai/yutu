// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

const (
	credFile   = "client_secret.json"
	tokenFile  = "youtube.token.json"
	credential = `{"client_id":"test"}`
	cacheToken = `{"access_token":"test"}`
	credB64    = "eyJjbGllbnRfaWQiOiJ0ZXN0In0="
	tokenB64   = "eyJhY2Nlc3NfdG9rZW4iOiJ0ZXN0In0="
)

func TestNewY2BService(t *testing.T) {
	rootDir, _ := os.Getwd()
	absCred := filepath.Join(rootDir, credFile)
	absToken := filepath.Join(rootDir, tokenFile)
	mockFS := fstest.MapFS{
		rootDir:   &fstest.MapFile{Mode: fs.ModeDir},
		credFile:  &fstest.MapFile{Data: []byte(credential)},
		tokenFile: &fstest.MapFile{Data: []byte(cacheToken)},
	}
	reader := strings.NewReader("in")
	writer := &bytes.Buffer{}

	type args struct {
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want Svc
	}{
		{
			name: "with all options - json",
			args: args{
				opts: []Option{
					WithCredential(credential, mockFS),
					WithCacheToken(cacheToken, mockFS),
					WithIO(reader, writer),
				},
			},
			want: &svc{
				Credential: credential,
				CacheToken: cacheToken,
				credFile:   credFile,
				ctx:        context.Background(),
				in:         reader,
				out:        writer,
			},
		},
		{
			name: "with all options - base64",
			args: args{
				opts: []Option{
					WithCredential(credB64, mockFS),
					WithCacheToken(tokenB64, mockFS),
					WithIO(reader, writer),
				},
			},
			want: &svc{
				Credential: credential,
				CacheToken: cacheToken,
				credFile:   credFile,
				ctx:        context.Background(),
				in:         reader,
				out:        writer,
			},
		},
		{
			name: "with all options - file",
			args: args{
				opts: []Option{
					WithCredential(absCred, mockFS),
					WithCacheToken(absToken, mockFS),
					WithIO(reader, writer),
				},
			},
			want: &svc{
				Credential: credential,
				CacheToken: cacheToken,
				credFile:   absCred,
				tokenFile:  tokenFile,
				ctx:        context.Background(),
				in:         reader,
				out:        writer,
			},
		},
		{
			name: "with no options",
			args: args{
				opts: []Option{},
			},
			want: &svc{
				credFile: credFile,
				ctx:      context.Background(),
				in:       os.Stdin,
				out:      os.Stdout,
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewY2BService(tt.args.opts...).(*svc)
				got.state = ""
				want := tt.want.(*svc)
				want.state = ""
				if !reflect.DeepEqual(got, want) {
					t.Errorf("NewY2BService() = %v, want %v", got, want)
				}
			},
		)
	}
}

func TestGetService_Errors(t *testing.T) {
	fsys := fstest.MapFS{}
	s := NewY2BService(WithCredential("not-a-valid-json-or-b64", fsys)).(*svc)

	_, err := s.GetService()
	if err == nil {
		t.Fatalf("GetService() error = nil, want non-nil")
	}
	if !strings.Contains(
		err.Error(), parseSecretFailed,
	) && !strings.Contains(err.Error(), "failed to parse client secret") {
		t.Fatalf("GetService() error = %v, want contains %q", err, parseSecretFailed)
	}
}
