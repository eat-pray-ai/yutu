// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/spf13/cobra"
)

const (
	authShort  = "Authenticate with YouTube API"
	authLong   = "Authenticate with YouTube API to access and manage YouTube resources."
	credUsage  = "Path to client secret file, or base64 encoded string, or json string"
	cacheUsage = "Path to token cache file, or base64 encoded string, or json string"
)

var (
	credential string
	cacheToken string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: authShort,
	Long:  authLong,
	Run: func(cmd *cobra.Command, args []string) {
		auth.NewY2BService(
			auth.WithCredential(credential, pkg.Root.FS()),
			auth.WithCacheToken(cacheToken, pkg.Root.FS()),
		).GetService()
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
}
