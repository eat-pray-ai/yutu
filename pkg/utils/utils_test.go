// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestGetFileName(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with extension",
			args: args{file: "example.txt"},
			want: "example",
		},
		{
			name: "without extension",
			args: args{file: "example"},
			want: "example",
		},
		{
			name: "with path",
			args: args{file: "/foo/bar/example.txt"},
			want: "example",
		},
		{
			name: "double extension",
			args: args{file: "archive.tar.gz"},
			want: "archive.tar",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := GetFileName(tt.args.file); got != tt.want {
					t.Errorf("GetFileName() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestIsJson(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid json",
			args: args{s: `{"key": "value"}`},
			want: true,
		},
		{
			name: "invalid json",
			args: args{s: `{"key": "value"`},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := IsJson(tt.args.s); got != tt.want {
					t.Errorf("IsJson() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestIsInteractive(t *testing.T) {
	devNull, _ := os.Open(os.DevNull)
	defer func() { _ = devNull.Close() }()

	tests := []struct {
		name   string
		writer io.Writer
		ci     string
		want   bool
	}{
		{
			name:   "bytes.Buffer is non-interactive",
			writer: &bytes.Buffer{},
			want:   false,
		},
		{
			name:   "os.File without TTY is non-interactive",
			writer: devNull,
			want:   false,
		},
		{
			name:   "CI env forces non-interactive",
			writer: os.Stdout,
			ci:     "true",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if tt.ci != "" {
					t.Setenv("CI", tt.ci)
				}
				if got := IsInteractive(tt.writer); got != tt.want {
					t.Errorf("IsInteractive() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestPrintJSON(t *testing.T) {
	tests := []struct {
		name        string
		data        any
		interactive bool
		wantWriter  string
	}{
		{
			name:        "interactive indented json",
			data:        map[string]string{"key": "value"},
			interactive: true,
			wantWriter:  "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name:        "non-interactive compact json",
			data:        map[string]string{"key": "value"},
			interactive: false,
			wantWriter:  "{\"key\":\"value\"}\n",
		},
		{
			name:        "nil data interactive",
			data:        nil,
			interactive: true,
			wantWriter:  "null\n",
		},
		{
			name:        "nil data non-interactive",
			data:        nil,
			interactive: false,
			wantWriter:  "null\n",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				orig := IsInteractive
				IsInteractive = func(any) bool { return tt.interactive }
				defer func() { IsInteractive = orig }()

				writer := &bytes.Buffer{}
				if err := PrintJSON(tt.data, writer); err != nil {
					t.Fatalf("PrintJSON() error = %v", err)
				}
				if gotWriter := writer.String(); gotWriter != tt.wantWriter {
					t.Errorf("PrintJSON() = %v, want %v", gotWriter, tt.wantWriter)
				}
			},
		)
	}
}

func TestPrintYAML(t *testing.T) {
	tests := []struct {
		name       string
		data       any
		wantWriter string
	}{
		{
			name:       "simple yaml",
			data:       map[string]string{"key": "value"},
			wantWriter: "key: value\n\n",
		},
		{
			name:       "nil data",
			data:       nil,
			wantWriter: "null\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				writer := &bytes.Buffer{}
				if err := PrintYAML(tt.data, writer); err != nil {
					t.Fatalf("PrintYAML() error = %v", err)
				}
				if gotWriter := writer.String(); gotWriter != tt.wantWriter {
					t.Errorf("PrintYAML() = %v, want %v", gotWriter, tt.wantWriter)
				}
			},
		)
	}
}

func TestPrintJSON_Error(t *testing.T) {
	if err := PrintJSON(func() {}, &bytes.Buffer{}); err == nil {
		t.Fatal("expected marshal error, got nil")
	}
}

func TestPrintYAML_Error(t *testing.T) {
	data := func() {}
	if err := PrintYAML(data, &bytes.Buffer{}); err == nil {
		t.Fatal("expected marshal error, got nil")
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestPrintJSON_WriteError(t *testing.T) {
	if err := PrintJSON(map[string]string{"key": "value"}, errWriter{}); err == nil {
		t.Fatal("expected write error, got nil")
	}
}

func TestPrintYAML_WriteError(t *testing.T) {
	if err := PrintYAML(map[string]string{"key": "value"}, errWriter{}); err == nil {
		t.Fatal("expected write error, got nil")
	}
}

func TestResetFlags(t *testing.T) {
	flags := pflag.NewFlagSet("reset", pflag.ContinueOnError)
	omittedBool, explicitBool, absentBool := new(true), new(true), new(true)
	flags.BoolVar(omittedBool, "omitted-bool", true, "")
	flags.BoolVar(explicitBool, "explicit-bool", true, "")
	var omittedTags, explicitTags []string
	absentTags := []string{"keep"}
	flags.StringSliceVar(&omittedTags, "omitted-tags", []string{"default"}, "")
	flags.StringSliceVar(&explicitTags, "explicit-tags", []string{"default"}, "")
	if err := flags.Parse([]string{"--explicit-bool=false", "--explicit-tags="}); err != nil {
		t.Fatal(err)
	}

	ResetFlags(map[string]**bool{
		"omitted-bool":  &omittedBool,
		"explicit-bool": &explicitBool,
		"absent-bool":   &absentBool,
	}, flags)
	ResetFlags(map[string]*[]string{
		"omitted-tags":  &omittedTags,
		"explicit-tags": &explicitTags,
		"absent-tags":   &absentTags,
	}, flags)

	if omittedBool != nil || omittedTags != nil {
		t.Errorf("omitted flags must reset to nil: bool=%v tags=%v", omittedBool, omittedTags)
	}
	if explicitBool == nil || *explicitBool {
		t.Errorf("explicit false must remain non-nil and false: %v", explicitBool)
	}
	if explicitTags == nil || len(explicitTags) != 0 {
		t.Errorf("explicit empty tags must remain non-nil and empty: %#v", explicitTags)
	}
	if absentBool == nil || !*absentBool || !reflect.DeepEqual(absentTags, []string{"keep"}) {
		t.Errorf("unregistered flags must remain unchanged: bool=%v tags=%v", absentBool, absentTags)
	}
}

func TestExtractHl(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid language uri with hl",
			args: args{uri: "i18n://language/zh-CN"},
			want: "zh-CN",
		},
		{
			name: "valid region uri with hl",
			args: args{uri: "i18n://region/zh-CN"},
			want: "zh-CN",
		},
		{
			name: "valid language uri without hl",
			args: args{uri: "i18n://language/"},
			want: "",
		},
		{
			name: "valid region uri without hl",
			args: args{uri: "i18n://region/"},
			want: "",
		},
		{
			name: "invalid uri",
			args: args{uri: "i18n://invalid/zh-CN"},
			want: "",
		},
		{
			name: "empty uri",
			args: args{uri: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := ExtractHl(tt.args.uri); got != tt.want {
					t.Errorf("ExtractHl() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestRandomStage(t *testing.T) {
	s1, err := RandomStage()
	if err != nil {
		t.Fatalf("RandomStage() error = %v", err)
	}
	s2, err := RandomStage()
	if err != nil {
		t.Fatalf("RandomStage() error = %v", err)
	}
	if s1 == "" {
		t.Error("RandomStage() returned empty string")
	}
	if s1 == s2 {
		t.Error("RandomStage() returned same value twice")
	}
}

func TestConfirmPreRun(t *testing.T) {
	tests := []struct {
		name        string
		yes         bool
		input       string
		interactive bool
		wantErr     error
	}{
		{
			name:    "yes flag bypasses prompt",
			yes:     true,
			wantErr: nil,
		},
		{
			name:        "non-interactive without yes flag",
			yes:         false,
			interactive: false,
			wantErr:     ErrNotConfirmed,
		},
		{
			name:        "interactive confirm with y",
			yes:         false,
			interactive: true,
			input:       "y\n",
			wantErr:     nil,
		},
		{
			name:        "interactive confirm with Y",
			yes:         false,
			interactive: true,
			input:       "Y\n",
			wantErr:     nil,
		},
		{
			name:        "interactive deny with n",
			yes:         false,
			interactive: true,
			input:       "n\n",
			wantErr:     ErrNotConfirmed,
		},
		{
			name:        "interactive deny with empty",
			yes:         false,
			interactive: true,
			input:       "\n",
			wantErr:     ErrNotConfirmed,
		},
		{
			name:        "interactive deny with EOF",
			yes:         false,
			interactive: true,
			input:       "",
			wantErr:     ErrNotConfirmed,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				orig := IsInteractive
				IsInteractive = func(any) bool { return tt.interactive }
				defer func() { IsInteractive = orig }()

				cmd := &cobra.Command{Use: "test"}
				cmd.Flags().Bool("yes", false, "")
				if tt.yes {
					_ = cmd.Flags().Set("yes", "true")
				}

				var errBuf bytes.Buffer
				cmd.SetErr(&errBuf)
				cmd.SetIn(strings.NewReader(tt.input))

				err := ConfirmPreRun(cmd, "Would do something")
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ConfirmPreRun() error = %v, want %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestJoinKeywords(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"nil", nil, ""},
		{"single", []string{"MCP"}, "MCP"},
		{"multi-word quoted", []string{"AI agents", "MCP"}, `"AI agents" MCP`},
		{"tab quoted", []string{"a\tb"}, "\"a\tb\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.in)
			if got := JoinKeywords(tt.in); got != tt.want {
				t.Errorf("JoinKeywords() = %q, want %q", got, tt.want)
			}
			if !slices.Equal(in, tt.in) {
				t.Errorf("input mutated: %v", tt.in)
			}
		})
	}
}
