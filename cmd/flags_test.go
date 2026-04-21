// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// newTestCmd creates a minimal cobra.Command with mutation flags registered.
func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	AddMutationFlags(cmd)
	return cmd
}

func TestAddMutationFlags(t *testing.T) {
	cmd := newTestCmd()

	dr := cmd.Flags().Lookup("dry-run")
	if dr == nil {
		t.Fatal("expected --dry-run flag to be registered")
	}
	if dr.DefValue != "false" {
		t.Fatalf("expected --dry-run default to be false, got %s", dr.DefValue)
	}

	yes := cmd.Flags().Lookup("yes")
	if yes == nil {
		t.Fatal("expected --yes flag to be registered")
	}
	if yes.DefValue != "false" {
		t.Fatalf("expected --yes default to be false, got %s", yes.DefValue)
	}
}

func TestDryRun_DefaultFalse(t *testing.T) {
	cmd := newTestCmd()
	if DryRun(cmd) {
		t.Fatal("expected DryRun to be false by default")
	}
}

func TestDryRun_SetTrue(t *testing.T) {
	cmd := newTestCmd()
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatal(err)
	}
	if !DryRun(cmd) {
		t.Fatal("expected DryRun to be true after setting flag")
	}
}

func TestConfirm_DryRun(t *testing.T) {
	cmd := newTestCmd()
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	err := Confirm(cmd, "would delete %s", "video123")
	if !errors.Is(err, ErrDryRun) {
		t.Fatalf("expected ErrDryRun, got %v", err)
	}

	got := stdout.String()
	if !strings.Contains(got, "would delete video123") {
		t.Fatalf("expected message in stdout, got %q", got)
	}
}

func TestConfirm_Yes(t *testing.T) {
	cmd := newTestCmd()
	if err := cmd.Flags().Set("yes", "true"); err != nil {
		t.Fatal(err)
	}

	err := Confirm(cmd, "delete video?")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestConfirm_NonTTY(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetIn(strings.NewReader(""))

	// Use the internal helper with isTTY=false to simulate non-TTY stdin.
	err := confirmWithTTYOverride(cmd, false, "delete video?")
	if err == nil {
		t.Fatal("expected error for non-TTY stdin")
	}
	if !strings.Contains(
		err.Error(), "terminal required for confirmation; use --yes to skip",
	) {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestConfirm_TTY_Yes(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetIn(strings.NewReader("y\n"))

	var stderr bytes.Buffer
	cmd.SetErr(&stderr)

	err := confirmWithTTYOverride(cmd, true, "delete %s?", "video123")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	prompt := stderr.String()
	if !strings.Contains(prompt, "delete video123?") {
		t.Fatalf("expected prompt on stderr, got %q", prompt)
	}
	if !strings.Contains(prompt, "[y/N]") {
		t.Fatalf("expected [y/N] in prompt, got %q", prompt)
	}
}

func TestConfirm_TTY_No(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetIn(strings.NewReader("n\n"))

	var stderr bytes.Buffer
	cmd.SetErr(&stderr)

	err := confirmWithTTYOverride(cmd, true, "delete video?")
	if !errors.Is(err, errAborted) {
		t.Fatalf("expected errAborted, got %v", err)
	}
}

func TestConfirm_TTY_Default(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetIn(strings.NewReader("\n"))

	var stderr bytes.Buffer
	cmd.SetErr(&stderr)

	err := confirmWithTTYOverride(cmd, true, "delete video?")
	if !errors.Is(err, errAborted) {
		t.Fatalf("expected errAborted for empty input (default N), got %v", err)
	}
}

func TestConfirm_DryRunTakesPrecedence(t *testing.T) {
	cmd := newTestCmd()
	if err := cmd.Flags().Set("dry-run", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("yes", "true"); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	err := Confirm(cmd, "would delete video")
	if !errors.Is(err, ErrDryRun) {
		t.Fatalf("expected ErrDryRun even with --yes, got %v", err)
	}
}
