package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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

var Server = mcp.NewServer(
	&mcp.Implementation{Name: "yutu", Version: Version},
	&mcp.ServerOptions{
		PageSize:     99,
		KeepAlive:    13 * time.Second,
		HasResources: true,
		HasTools:     true,
	},
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: mcpShort,
	Long:  mcpLong,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		ctx := context.Background()
		addr := fmt.Sprintf(":%d", port)
		slog.InfoContext(
			ctx, "starting MCP server",
			"mode", mode,
			"version", Version,
		)

		switch mode {
		case "stdio":
			t := &mcp.LoggingTransport{
				Transport: &mcp.StdioTransport{},
				Writer:    os.Stderr,
			}
			err = Server.Run(ctx, t)
		case "http":
			handler := mcp.NewStreamableHTTPHandler(
				func(*http.Request) *mcp.Server {
					return Server
				}, nil,
			)
			slog.InfoContext(
				ctx, "http server configuration",
				"url", fmt.Sprintf("http://localhost:%d/mcp", port),
			)
			err = http.ListenAndServe(addr, handler)
		default:
			slog.ErrorContext(
				ctx, "invalid mode", "mode", mode, "valid_modes", "stdio, http",
			)
			os.Exit(1)
		}

		if err != nil {
			slog.ErrorContext(
				ctx, "starting server failed", "error", err, "mode", mode,
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
