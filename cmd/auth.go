// Copyright 2025 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"log/slog"
	"os"

	"github.com/eat-pray-ai/yutu/pkg"
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/eat-pray-ai/yutu/pkg/common"
	"github.com/spf13/cobra"
)

const (
	authShort      = "Authenticate with YouTube APIs"
	authLong       = "Authenticate with YouTube APIs to access and manage YouTube resources."
	credUsage      = "Path to client secret file, or Base64 encoded string, or JSON string (env: YUTU_CREDENTIAL)"
	cacheUsage     = "Path to token cache file, or Base64 encoded string, or JSON string (env: YUTU_CACHE_TOKEN)"
	noBrowserUsage = "Do not open browser automatically"
)

func init() {
	RootCmd.AddCommand(newAuthCmd())
}

type authOpts struct {
	credential string
	cacheToken string
	noBrowser  bool
}

func (o *authOpts) bindFlags(c *cobra.Command) {
	c.Flags().StringVarP(
		&o.credential, "credential", "c", "client_secret.json", credUsage,
	)
	c.Flags().StringVarP(
		&o.cacheToken, "cacheToken", "t", "youtube.token.json", cacheUsage,
	)
	c.Flags().BoolVarP(&o.noBrowser, "noBrowser", "B", false, noBrowserUsage)
}

func newAuthCmd() *cobra.Command {
	var opts authOpts

	c := &cobra.Command{
		Use:   "auth",
		Short: authShort,
		Long:  authLong,
		Run: func(c *cobra.Command, _ []string) {
			if _, err := auth.NewY2BService(
				auth.WithContext(c.Context()),
				auth.WithCredential(opts.credential, pkg.Root.FS()),
				auth.WithCacheToken(opts.cacheToken, pkg.Root.FS()),
				auth.WithRedirectURL(common.DefaultRedirectURL),
				auth.WithNoBrowser(opts.noBrowser),
			).GetService(); err != nil {
				slog.Error("authentication failed", "error", err)
				os.Exit(1)
			}
		},
	}

	opts.bindFlags(c)
	return c
}
