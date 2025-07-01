package utils

import (
	"bytes"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"reflect"
	"testing"
)

func TestBoolPtr(t *testing.T) {
	type args struct {
		b string
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "true",
			args: args{b: "true"},
			want: func() *bool {
				b := true
				return &b
			}(),
		},
		{
			name: "false",
			args: args{b: "false"},
			want: func() *bool {
				b := false
				return &b
			}(),
		},
		{
			name: "empty",
			args: args{b: ""},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := BoolPtr(tt.args.b); !reflect.DeepEqual(got, tt.want) {
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
		data  interface{}
		jpath string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name:       "empty jsonpath",
			args:       args{data: map[string]string{"key": "value"}, jpath: ""},
			wantWriter: "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name:       "simple json",
			args:       args{data: map[string]string{"key": "value"}, jpath: "$"},
			wantWriter: "[\n  {\n    \"key\": \"value\"\n  }\n]\n",
		},
		{
			name: "invalid jsonpath",
			args: args{
				data: map[string]string{"key": "value"}, jpath: "//",
			},
			wantWriter: "Invalid JSONPath: //\n",
		},
		{
			name: "extract specific field",
			args: args{
				data: map[string]interface{}{
					"key":     "value",
					"another": "field",
				}, jpath: "$.key",
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
				jpath: "$.*.key1",
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
				jpath: "$[*].key1",
			},
			wantWriter: "[\n  \"value1\",\n  \"value3\"\n]\n",
		},
		{
			name:       "nil data",
			args:       args{data: nil, jpath: "$"},
			wantWriter: "[\n  null\n]\n",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				writer := &bytes.Buffer{}
				PrintJSON(tt.args.data, tt.args.jpath, writer)
				if gotWriter := writer.String(); gotWriter != tt.wantWriter {
					t.Errorf("PrintJSON() = %v, want %v", gotWriter, tt.wantWriter)
				}
			},
		)
	}
}

func TestPrintYAML(t *testing.T) {
	type args struct {
		data  interface{}
		jpath string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name:       "empty jsonpath",
			args:       args{data: map[string]string{"key": "value"}, jpath: ""},
			wantWriter: "key: value\n\n",
		},
		{
			name:       "simple yaml",
			args:       args{data: map[string]string{"key": "value"}, jpath: "$"},
			wantWriter: "- key: value\n\n",
		},
		{
			name: "invalid jsonpath",
			args: args{
				data: map[string]string{"key": "value"}, jpath: "//",
			},
			wantWriter: "Invalid JSONPath: //\n",
		},
		{
			name: "extract specific field",
			args: args{
				data: map[string]interface{}{
					"key":     "value",
					"another": "field",
				}, jpath: "$.key",
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
				jpath: "$.*.key1",
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
				jpath: "$[*].key1",
			},
			wantWriter: "- value1\n- value3\n\n",
		},
		{
			name:       "nil data",
			args:       args{data: nil, jpath: "$"},
			wantWriter: "- null\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				writer := &bytes.Buffer{}
				PrintYAML(tt.args.data, tt.args.jpath, writer)
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
	b := BoolPtr("true")
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
