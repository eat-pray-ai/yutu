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
	Root, err = os.OpenRoot("/")
	if err != nil {
		panic(err)
	}

	Logger = slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	slog.SetDefault(Logger)
}
