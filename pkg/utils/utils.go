// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bufio"
	"cmp"
	"crypto/rand"
	"encoding/base64"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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

func PrintJSON(data any, writer io.Writer) (err error) {
	var marshalled []byte
	if IsInteractive(writer) {
		marshalled, err = json.Marshal(data, jsontext.WithIndent("  "))
	} else {
		marshalled, err = json.Marshal(data)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	_, err = fmt.Fprintln(writer, string(marshalled))
	return err
}

func PrintYAML(data any, writer io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// yaml.v3 may panic for unsupported values (for example, functions).
			// Preserve this as an ordinary error for CLI callers.
			err = fmt.Errorf("failed to marshal YAML: %v", r)
		}
	}()
	marshalled, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	_, err = fmt.Fprintln(writer, string(marshalled))
	return err
}

func RandomStage() (string, error) {
	b := make([]byte, 128)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(b)
	return state, nil
}

func GetFileName(file string) string {
	base := filepath.Base(file)
	fileName := base[:len(base)-len(filepath.Ext(base))]
	return fileName
}

func IsJson(s string) bool {
	var js jsontext.Value
	return json.Unmarshal([]byte(s), &js) == nil
}

func ResetFlags[T any](values map[string]*T, flagSet *pflag.FlagSet) {
	for name, value := range values {
		flag := flagSet.Lookup(name)
		if flag != nil && !flag.Changed {
			var zero T
			*value = zero
		}
	}
}

func ResolveCredential(cmd *cobra.Command, flag, env, value string) (string, error) {
	if !cmd.Flags().Changed(flag) {
		value = cmp.Or(os.Getenv(env), value)
	}
	if strings.TrimSpace(value) == "" {
		return "", nil
	}
	// Keep inline values intact; make paths independent of the agent's cwd.
	if jsontext.Value(value).IsValid() {
		return value, nil
	}
	if _, err := base64.StdEncoding.DecodeString(value); err == nil {
		return value, nil
	}
	return filepath.Abs(value)
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
