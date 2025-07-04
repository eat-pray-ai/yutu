package cmd

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"log"
	"time"
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
	server.WithLogging(),
	server.WithRecovery(),
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: mcpShort,
	Long:  mcpLong,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		interval := 13 * time.Second
		addr := fmt.Sprintf(":%d", port)
		message := fmt.Sprintf("%s server listening on %s", mode, addr)

		switch mode {
		case "stdio":
			err = server.ServeStdio(MCP)
		case "http":
			httpServer := server.NewStreamableHTTPServer(
				MCP,
				server.WithHeartbeatInterval(interval),
			)
			log.Printf("%s/%s\n", message, "mcp")
			err = httpServer.Start(addr)
		}

		if err != nil {
			log.Printf("Server error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().StringVarP(&mode, "mode", "m", "stdio", modeUsage)
	mcpCmd.Flags().IntVarP(&port, "port", "p", 8216, portUsage)
}
