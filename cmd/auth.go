// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/spf13/cobra"
)

const (
	authShort  = "Authenticate with YouTube APIs"
	authLong   = "Authenticate with YouTube APIs to access and manage YouTube resources."
	credUsage  = "Path to client secret file, or base64 encoded string, or json string (env: YUTU_CREDENTIAL)"
	cacheUsage = "Path to token cache file, or base64 encoded string, or json string (env: YUTU_CACHE_TOKEN)"
)

var (
	credential string
	cacheToken string
	authPort   int
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: authShort,
	Long:  authLong,
	Run: func(cmd *cobra.Command, args []string) {
		redirectURL := fmt.Sprintf("http://localhost:%d", authPort)
		if _, err := auth.NewY2BService(
			auth.WithCredential(credential, pkg.Root.FS()),
			auth.WithCacheToken(cacheToken, pkg.Root.FS()),
			auth.WithRedirectURL(redirectURL),
		).GetService(); err != nil {
			slog.Error("authentication failed", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(
		&credential, "credential", "c", "client_secret.json", credUsage,
	)
	authCmd.Flags().StringVarP(
		&cacheToken, "cacheToken", "t", "youtube.token.json", cacheUsage,
	)
	authCmd.Flags().IntVarP(
		&authPort, "port", "p", 8216, "Port for OAuth redirect URL",
	)
}
