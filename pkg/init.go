// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log/slog"
	"os"
	"runtime"
)

var (
	RootDir *string
	Root    *os.Root
	logger  *slog.Logger
)

func init() {
	if logger == nil {
		initLogger()
	}

	if RootDir == nil {
		initRootDir()
	}
}

func initLogger() {
	logLevel := slog.LevelInfo
	if lvl, ok := os.LookupEnv("YUTU_LOG_LEVEL"); ok {
		switch lvl {
		case "DEBUG", "debug":
			logLevel = slog.LevelDebug
		case "INFO", "info":
			logLevel = slog.LevelInfo
		case "WARN", "warn":
			logLevel = slog.LevelWarn
		case "ERROR", "error":
			logLevel = slog.LevelError
		}
	}
	logger = slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			},
		),
	)

	slog.SetDefault(logger)
}

func initRootDir() {
	var err error
	rootDir, ok := os.LookupEnv("YUTU_ROOT")
	if !ok {
		rootDir, err = os.Getwd()
		if err != nil {
			rootDir = "/"
			if runtime.GOOS == "windows" {
				rootDir = os.Getenv("SystemDrive") + `\`
			}
			slog.Debug(getWdFailed, "error", err, "fallback", rootDir)
		}
	}
	RootDir = &rootDir
	Root, err = os.OpenRoot(*RootDir)
	if err != nil {
		slog.Error(openRootFailed, "dir", *RootDir, "error", err)
		panic(err)
	}
	slog.Debug("Root directory set", "dir", *RootDir)
}
