// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"log/slog"
	"os"
	"runtime"
)

const (
	PartsUsage  = "Comma separated parts"
	MRUsage     = "The maximum number of items that should be returned, 0 for no limit"
	TableUsage  = "json, yaml, or table"
	SilentUsage = "json, yaml, or silent"
	JPUsage     = "JSONPath expression to filter the output"
	JsonMIME    = "application/json"
	PerPage     = 20

	getWdFailed    = "failed to get working directory"
	openRootFailed = "failed to open root directory"
)

var (
	RootDir *string
	Root    *os.Root
	Logger  *slog.Logger
)

func init() {
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
	Logger = slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			},
		),
	)

	slog.SetDefault(Logger)

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
