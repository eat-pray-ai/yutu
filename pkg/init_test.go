package pkg

import (
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
		{"default", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.level != "" {
				t.Setenv("YUTU_LOG_LEVEL", tt.level)
			} else {
				os.Unsetenv("YUTU_LOG_LEVEL")
			}
			initLogger()
			if logger == nil {
				t.Error("logger is nil")
			}
		})
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

	t.Run("with env var", func(t *testing.T) {
		wd, _ := os.Getwd()
		t.Setenv("YUTU_ROOT", wd)
		initRootDir()
		if *RootDir != wd {
			t.Errorf("expected %s, got %s", wd, *RootDir)
		}
		if Root == nil {
			t.Error("Root is nil")
		}
	})

	t.Run("without env var", func(t *testing.T) {
		os.Unsetenv("YUTU_ROOT")
		initRootDir()
		if RootDir == nil {
			t.Error("RootDir is nil")
		}
		// Should fallback to CWD
		wd, _ := os.Getwd()
		if *RootDir != wd {
			t.Errorf("expected %s, got %s", wd, *RootDir)
		}
	})
}
