// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"

	"gopkg.in/yaml.v3"
)

var IsInteractive = func(v any) bool {
	if os.Getenv("CI") != "" {
		return false
	}
	if f, ok := v.(*os.File); ok {
		return term.IsTerminal(int(f.Fd()))
	}
	return false
}

func PrintJSON(data any, writer io.Writer) {
	var marshalled []byte
	if IsInteractive(writer) {
		marshalled, _ = json.MarshalIndent(data, "", "  ")
	} else {
		marshalled, _ = json.Marshal(data)
	}
	_, _ = fmt.Fprintln(writer, string(marshalled))
}

func PrintYAML(data any, writer io.Writer) {
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

func StrToBoolPtr(b *string) *bool {
	if b == nil || *b == "" || strings.ToLower(strings.TrimSpace(*b)) == "null" {
		return nil
	}
	return new(*b == "true")
}

func BoolToStrPtr(b *bool) *string {
	if b == nil {
		return new("")
	}

	return new(strconv.FormatBool(*b))
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

var ErrNotConfirmed = errors.New("operation not confirmed (pass confirm: true or rerun with --yes)")

func ConfirmPreRun(cmd *cobra.Command, msg string) error {
	if yes, _ := cmd.Flags().GetBool("yes"); yes {
		return nil
	}

	if !IsInteractive(cmd.InOrStdin()) {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), msg)
		return ErrNotConfirmed
	}

	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s [y/N] ", msg)

	scanner := bufio.NewScanner(cmd.InOrStdin())
	if !scanner.Scan() {
		return ErrNotConfirmed
	}
	answer := strings.TrimSpace(scanner.Text())
	if answer != "y" && answer != "Y" {
		return ErrNotConfirmed
	}
	return nil
}

func HandleCmdError(err error, cmd *cobra.Command) {
	if err == nil {
		return
	}
	_ = cmd.Help()
	cmd.PrintErrf("Error: %v\n", err)
}
