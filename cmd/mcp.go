package cmd

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	modeUsage = "stdio, http or sse"
	portUsage = "Port to listen on for HTTP or SSE mode"
)

var (
	mode string
	port int
)

var MCP = server.NewMCPServer(
	"yutu", Version,
	server.WithToolCapabilities(true),
	server.WithRecovery(),
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start mcp server",
	Long:  "Start mcp server to handle requests from clients",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		interval := 13 * time.Second
		addr := fmt.Sprintf(":%d", port)
		baseURL := fmt.Sprintf("http://localhost:%d", port)
		message := fmt.Sprintf("Starting MCP server: %s", baseURL)

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
		case "sse":
			sse := server.NewSSEServer(
				MCP, server.WithBaseURL(baseURL),
				server.WithKeepAlive(true),
				server.WithKeepAliveInterval(interval),
			)
			log.Printf("%s/%s\n", message, "sse")
			err = sse.Start(addr)
		}

		if err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(mcpCmd)

	mcpCmd.Flags().StringVarP(&mode, "mode", "m", "stdio", modeUsage)
	mcpCmd.Flags().IntVarP(&port, "port", "p", 8080, portUsage)
}
