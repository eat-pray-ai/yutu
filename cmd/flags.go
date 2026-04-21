// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eat-pray-ai/yutu/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// ErrDryRun is returned by Confirm when --dry-run is set.
var ErrDryRun = utils.ErrDryRun

// errAborted is returned by Confirm when the user declines the prompt.
var errAborted = errors.New("aborted")

// AddMutationFlags registers --dry-run and --yes on the given command.
func AddMutationFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(
		"dry-run", false,
		"Print what would happen without calling the API",
	)
	cmd.Flags().Bool(
		"yes", false,
		"Skip the confirmation prompt",
	)
}

// DryRun returns whether --dry-run was passed.
func DryRun(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("dry-run")
	return v
}

// Confirm checks dry-run and yes flags, then prompts the user for
// confirmation if running in an interactive terminal.
//
// Returns ErrDryRun if --dry-run is set, nil if --yes is set or the user
// confirms, and errAborted if the user declines.
func Confirm(cmd *cobra.Command, format string, args ...any) error {
	isTTY := term.IsTerminal(int(os.Stdin.Fd()))
	return confirmWithTTYOverride(cmd, isTTY, format, args...)
}

// confirmWithTTYOverride is an internal helper that accepts a pre-computed
// isTTY value so tests can exercise the TTY and non-TTY paths without
// needing a real terminal.
func confirmWithTTYOverride(cmd *cobra.Command, isTTY bool, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)

	// 1. --dry-run takes precedence over everything.
	if DryRun(cmd) {
		fmt.Fprintln(cmd.OutOrStdout(), msg)
		return ErrDryRun
	}

	// 2. --yes skips the prompt entirely.
	yes, _ := cmd.Flags().GetBool("yes")
	if yes {
		return nil
	}

	// 3. Non-interactive stdin requires --yes.
	if !isTTY {
		return errors.New("terminal required for confirmation; use --yes to skip")
	}

	// 4. Interactive prompt.
	_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "%s [y/N] ", msg)

	scanner := bufio.NewScanner(cmd.InOrStdin())
	if !scanner.Scan() {
		return errAborted
	}
	answer := strings.TrimSpace(scanner.Text())

	if answer == "y" || answer == "Y" {
		return nil
	}
	return errAborted
}
