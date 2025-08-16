package cmd

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

var Logger *slog.Logger

func init() {
	Logger = slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	slog.SetDefault(Logger)
}
