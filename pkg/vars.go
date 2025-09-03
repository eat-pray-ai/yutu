package pkg

import (
	"log/slog"
	"os"
)

const (
	PartsUsage  = "Comma separated parts"
	MRUsage     = "The maximum number of items that should be returned"
	TableUsage  = "json, yaml, or table"
	SilentUsage = "json, yaml, or silent"
	JPUsage     = "JSONPath expression to filter the output"
	JsonMIME    = "application/json"
)

var (
	Root   *os.Root
	Logger *slog.Logger
)

func init() {
	var err error
	rootDir, ok := os.LookupEnv("YUTU_ROOT")
	if !ok {
		rootDir = "/"
	}
	Root, err = os.OpenRoot(rootDir)
	if err != nil {
		panic(err)
	}

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
}
