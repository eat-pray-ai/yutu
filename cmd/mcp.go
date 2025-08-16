package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

const (
	mcpShort  = "Start MCP server"
	mcpLong   = "Start MCP server to handle requests from clients"
	modeUsage = "stdio, or http"
	portUsage = "Port to listen on for HTTP or SSE mode"
)

var (
	mode string
	port int
)

var MCP = server.NewMCPServer(
	"yutu", Version,
	server.WithToolCapabilities(true),
	server.WithResourceCapabilities(true, true),
	server.WithLogging(),
	server.WithRecovery(),
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: mcpShort,
	Long:  mcpLong,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		var err error
		interval := 13 * time.Second
		addr := fmt.Sprintf(":%d", port)
		slog.InfoContext(
			ctx, "starting MCP server",
			"mode", mode,
			"version", Version,
		)

		switch mode {
		case "stdio":
			err = server.ServeStdio(MCP)
		case "http":
			httpServer := server.NewStreamableHTTPServer(
				MCP,
				server.WithHeartbeatInterval(interval),
			)

			slog.InfoContext(
				ctx, "http server configuration",
				"url", fmt.Sprintf("http://localhost:%d/mcp", port),
				"heartbeat_interval", interval,
			)
			err = httpServer.Start(addr)
		default:
			slog.ErrorContext(
				ctx, "invalid server mode",
				"mode", mode,
				"valid_modes", []string{"stdio", "http"},
			)
			os.Exit(1)
		}

		if err != nil {
			slog.ErrorContext(
				ctx, "starting server failed",
				"error", err,
				"mode", mode,
			)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().StringVarP(&mode, "mode", "m", "stdio", modeUsage)
	mcpCmd.Flags().IntVarP(&port, "port", "p", 8216, portUsage)
}
