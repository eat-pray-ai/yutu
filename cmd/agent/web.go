// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package agent

import (
	"slices"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start agent in web mode",
	Long:  "start agent in web mode with additional sub-servers specified by sublaunchers",
	Run: func(cmd *cobra.Command, args []string) {
		launcherArgs := []string{"web"}

		launcherArgs = append(launcherArgs, "-idle-timeout", idleTimeout.String())
		launcherArgs = append(launcherArgs, "-port", strconv.Itoa(port))
		launcherArgs = append(launcherArgs, "-read-timeout", readTimeout.String())
		launcherArgs = append(launcherArgs, "-write-timeout", writeTimeout.String())

		if slices.Contains(sublaunchers, "api") {
			launcherArgs = append(
				launcherArgs, "api",
				"-sse-write-timeout", sseWriteTimeout.String(),
				"-webui_address", webuiAddress,
			)
		}
		if slices.Contains(sublaunchers, "a2a") {
			launcherArgs = append(launcherArgs, "a2a", "-a2a_agent_url", a2aAgentURL)
		}
		if slices.Contains(sublaunchers, "webui") {
			launcherArgs = append(launcherArgs, "webui", "-api_server_address", apiServerAddress)
		}

		launch(cmd.Context(), launcherArgs)
	},
}

func init() {
	agentCmd.AddCommand(webCmd)

	// web flags
	webCmd.Flags().DurationVarP(&idleTimeout, "idleTimeout", "i", time.Minute, "Server idle timeout")
	webCmd.Flags().IntVarP(&port, "port", "p", 8080, "Localhost port for the server")
	webCmd.Flags().DurationVarP(&readTimeout, "readTimeout", "r", 15*time.Second, "Server read timeout")
	webCmd.Flags().DurationVarP(&writeTimeout, "writeTimeout", "w", 15*time.Second, "Server write timeout")
	webCmd.Flags().StringSliceVarP(&sublaunchers, "sublaunchers", "l", []string{"api"}, "One or more sublaunchers: api, a2a, and webui")

	// api flags
	webCmd.Flags().DurationVarP(&sseWriteTimeout, "sseWriteTimeout", "s", 2*time.Minute, "SSE server write timeout")
	webCmd.Flags().StringVarP(&webuiAddress, "webuiAddress", "W", "localhost:8080", "ADK WebUI address")

	// a2a flags
	webCmd.Flags().StringVarP(&a2aAgentURL, "a2aAgentUrl", "u", "http://localhost:8080", "A2A host URL")

	// webui flags
	webCmd.Flags().StringVarP(&apiServerAddress, "apiServerAddress", "U", "http://localhost:8080/api", "ADK REST API server address")
}
