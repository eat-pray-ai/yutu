// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ohler55/ojg/jp"
	"github.com/spf13/pflag"

	"gopkg.in/yaml.v3"
)

func PrintJSON(data interface{}, jsonpath string, writer io.Writer) {
	j, err := jp.ParseString(jsonpath)
	if err != nil && jsonpath != "" {
		_, _ = fmt.Fprintln(writer, "Invalid JSONPath:", jsonpath)
		return
	} else if jsonpath != "" {
		data = j.Get(data)
	}
	marshalled, _ := json.MarshalIndent(data, "", "  ")
	_, _ = fmt.Fprintln(writer, string(marshalled))
}

func PrintYAML(data interface{}, jsonpath string, writer io.Writer) {
	j, err := jp.ParseString(jsonpath)
	if err != nil && jsonpath != "" {
		_, _ = fmt.Fprintln(writer, "Invalid JSONPath:", jsonpath)
		return
	} else if jsonpath != "" {
		data = j.Get(data)
	}
	marshalled, _ := yaml.Marshal(data)
	_, _ = fmt.Fprintln(writer, string(marshalled))
}

func OpenURL(url string) error {
	var err error
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

func RandomStage() string {
	b := make([]byte, 128)
	_, _ = rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func GetFileName(file string) string {
	base := filepath.Base(file)
	fileName := base[:len(base)-len(filepath.Ext(base))]
	return fileName
}

func IsJson(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func BoolPtr(b string) *bool {
	if b == "" || strings.ToLower(strings.TrimSpace(b)) == "null" {
		return nil
	}
	val := b == "true"
	return &val
}

func ResetBool(m map[string]**bool, flagSet *pflag.FlagSet) {
	for k := range m {
		flag := flagSet.Lookup(k)
		if flag != nil && !flag.Changed {
			*m[k] = nil
		}
	}
}

func ExtractHl(uri string) string {
	pattern := `i18n://(?:language|region)/([^/]+)`
	matches := regexp.MustCompile(pattern).FindStringSubmatch(uri)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func Errf(format string, args ...any) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf(format, args...)}},
		IsError: true,
	}
}
