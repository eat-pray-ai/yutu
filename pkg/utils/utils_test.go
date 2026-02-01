// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestStrToBoolPtr(t *testing.T) {
	type args struct {
		b *string
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "true",
			args: args{b: jsonschema.Ptr("true")},
			want: func() *bool {
				b := true
				return &b
			}(),
		},
		{
			name: "false",
			args: args{b: jsonschema.Ptr("false")},
			want: func() *bool {
				b := false
				return &b
			}(),
		},
		{
			name: "empty",
			args: args{b: jsonschema.Ptr("")},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := StrToBoolPtr(tt.args.b); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BoolPtr() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

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

func TestPrintJSON(t *testing.T) {
	type args struct {
		data     interface{}
		jsonpath string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name:       "empty jsonpath",
			args:       args{data: map[string]string{"key": "value"}, jsonpath: ""},
			wantWriter: "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name:       "simple json",
			args:       args{data: map[string]string{"key": "value"}, jsonpath: "$"},
			wantWriter: "[\n  {\n    \"key\": \"value\"\n  }\n]\n",
		},
		{
			name: "invalid jsonpath",
			args: args{
				data: map[string]string{"key": "value"}, jsonpath: "//",
			},
			wantWriter: "Invalid JSONPath: //\n",
		},
		{
			name: "extract specific field",
			args: args{
				data: map[string]interface{}{
					"key":     "value",
					"another": "field",
				}, jsonpath: "$.key",
			},
			wantWriter: "[\n  \"value\"\n]\n",
		},
		{
			name: "nested jsonpath",
			args: args{
				data: map[string]interface{}{
					"item1": map[string]string{"key1": "value1"},
					"item2": map[string]string{"key2": "value2"},
					"count": 2,
				},
				jsonpath: "$.*.key1",
			},
			wantWriter: "[\n  \"value1\"\n]\n",
		},
		{
			name: "array jsonpath",
			args: args{
				data: []map[string]string{
					{"key1": "value1"},
					{"key2": "value2"},
					{"key1": "value3"},
				},
				jsonpath: "$[*].key1",
			},
			wantWriter: "[\n  \"value1\",\n  \"value3\"\n]\n",
		},
		{
			name:       "nil data",
			args:       args{data: nil, jsonpath: "$"},
			wantWriter: "[\n  null\n]\n",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				writer := &bytes.Buffer{}
				PrintJSON(tt.args.data, tt.args.jsonpath, writer)
				if gotWriter := writer.String(); gotWriter != tt.wantWriter {
					t.Errorf("PrintJSON() = %v, want %v", gotWriter, tt.wantWriter)
				}
			},
		)
	}
}

func TestPrintYAML(t *testing.T) {
	type args struct {
		data     interface{}
		jsonpath string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name:       "empty jsonpath",
			args:       args{data: map[string]string{"key": "value"}, jsonpath: ""},
			wantWriter: "key: value\n\n",
		},
		{
			name:       "simple yaml",
			args:       args{data: map[string]string{"key": "value"}, jsonpath: "$"},
			wantWriter: "- key: value\n\n",
		},
		{
			name: "invalid jsonpath",
			args: args{
				data: map[string]string{"key": "value"}, jsonpath: "//",
			},
			wantWriter: "Invalid JSONPath: //\n",
		},
		{
			name: "extract specific field",
			args: args{
				data: map[string]interface{}{
					"key":     "value",
					"another": "field",
				}, jsonpath: "$.key",
			},
			wantWriter: "- value\n\n",
		},
		{
			name: "nested jsonpath",
			args: args{
				data: map[string]interface{}{
					"item1": map[string]string{"key1": "value1"},
					"item2": map[string]string{"key2": "value2"},
					"count": 2,
				},
				jsonpath: "$.*.key1",
			},
			wantWriter: "- value1\n\n",
		},
		{
			name: "array jsonpath",
			args: args{
				data: []map[string]string{
					{"key1": "value1"},
					{"key2": "value2"},
					{"key1": "value3"},
				},
				jsonpath: "$[*].key1",
			},
			wantWriter: "- value1\n- value3\n\n",
		},
		{
			name:       "nil data",
			args:       args{data: nil, jsonpath: "$"},
			wantWriter: "- null\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				writer := &bytes.Buffer{}
				PrintYAML(tt.args.data, tt.args.jsonpath, writer)
				if gotWriter := writer.String(); gotWriter != tt.wantWriter {
					t.Errorf("PrintYAML() = %v, want %v", gotWriter, tt.wantWriter)
				}
			},
		)
	}
}

func TestResetBool(t *testing.T) {
	type args struct {
		m       map[string]**bool
		flagSet *pflag.FlagSet
	}
	b := jsonschema.Ptr(true)
	cmd := &cobra.Command{}
	cmd.Flags().BoolVar(b, "flag", false, "")
	tests := []struct {
		name string
		args args
	}{
		{
			name: "reset bool flags",
			args: args{
				m: map[string]**bool{
					"flag": &b,
				},
				flagSet: cmd.Flags(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ResetBool(tt.args.m, tt.args.flagSet)
				if b != nil {
					t.Errorf("ResetBool() = %v, want nil", *b)
				}
			},
		)
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

func TestBoolToStrPtr(t *testing.T) {
	bTrue := true
	bFalse := false
	type args struct {
		b *bool
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "true",
			args: args{b: &bTrue},
			want: jsonschema.Ptr("true"),
		},
		{
			name: "false",
			args: args{b: &bFalse},
			want: jsonschema.Ptr("false"),
		},
		{
			name: "nil",
			args: args{b: nil},
			want: jsonschema.Ptr(""),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := BoolToStrPtr(tt.args.b); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("BoolToStrPtr() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestHandleCmdError(t *testing.T) {
	tests := []struct {
		name    string
		input   error
		wantOut string
		wantErr string
	}{
		{
			name:    "with error",
			input:   fmt.Errorf("some error"),
			wantOut: "help called",
			wantErr: "Error: some error\n",
		},
		{
			name:    "without error",
			input:   nil,
			wantOut: "",
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				cmd := &cobra.Command{Use: "test"}
				var outBuf, errBuf bytes.Buffer
				cmd.SetOut(&outBuf)
				cmd.SetErr(&errBuf)
				cmd.SetHelpFunc(
					func(c *cobra.Command, args []string) {
						_, _ = fmt.Fprint(c.OutOrStdout(), "help called")
					},
				)

				HandleCmdError(tt.input, cmd)

				if gotOut := outBuf.String(); gotOut != tt.wantOut {
					t.Fatalf(
						"unexpected stdout output, got %q, want %q", gotOut, tt.wantOut,
					)
				}

				if gotErr := errBuf.String(); gotErr != tt.wantErr {
					t.Fatalf(
						"unexpected stderr output, got %q, want %q", gotErr, tt.wantErr,
					)
				}
			},
		)
	}
}
