// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"errors"
	"os"
	"testing"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"DEBUG uppercase", "DEBUG"},
		{"INFO uppercase", "INFO"},
		{"WARN uppercase", "WARN"},
		{"ERROR uppercase", "ERROR"},
		{"default", ""},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if tt.level != "" {
					t.Setenv("YUTU_LOG_LEVEL", tt.level)
				} else {
					_ = os.Unsetenv("YUTU_LOG_LEVEL")
				}
				initLogger()
				if logger == nil {
					t.Error("logger is nil")
				}
			},
		)
	}
}

func TestInitRootDir(t *testing.T) {
	// Save original RootDir and restore after test
	origRootDir := RootDir
	origRoot := Root
	defer func() {
		RootDir = origRootDir
		Root = origRoot
	}()

	t.Run(
		"with env var", func(t *testing.T) {
			wd, _ := os.Getwd()
			t.Setenv("YUTU_ROOT", wd)
			initRootDir()
			if *RootDir != wd {
				t.Errorf("expected %s, got %s", wd, *RootDir)
			}
			if Root == nil {
				t.Error("Root is nil")
			}
		},
	)

	t.Run(
		"without env var", func(t *testing.T) {
			_ = os.Unsetenv("YUTU_ROOT")
			initRootDir()
			if RootDir == nil {
				t.Error("RootDir is nil")
			}
			// Should fallback to CWD
			wd, _ := os.Getwd()
			if *RootDir != wd {
				t.Errorf("expected %s, got %s", wd, *RootDir)
			}
		},
	)
}

func TestOpenUploadFile(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		wantBlock bool
	}{
		{"json file", "client_secret.json", true},
		{"token file", "youtube.token.json", true},
		{"env file", ".env", true},
		{"key file", "id_rsa.key", true},
		{"pem file", "server.pem", true},
		{"crt file", "tls.crt", true},
		{"cert file", "client.cert", true},
		{"pfx file", "bundle.pfx", true},
		{"p12 file", "cert.p12", true},
		{"uppercase JSON", "SECRET.JSON", true},
		{"valid mp4", "video.mp4", false},
		{"valid srt", "caption.srt", false},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := OpenFile(tt.file)
				if tt.wantBlock && !errors.Is(err, ErrBlockedFile) {
					t.Errorf(
						"OpenUploadFile(%q) error = %v, want ErrBlockedFile", tt.file, err,
					)
				}
				if !tt.wantBlock && errors.Is(err, ErrBlockedFile) {
					t.Errorf("OpenUploadFile(%q) blocked unexpectedly", tt.file)
				}
			},
		)
	}
}
