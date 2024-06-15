package cmd

import (
	"github.com/eat-pray-ai/yutu/pkg/auth"
	"github.com/spf13/cobra"
)

var (
	credential string
	cacheToken string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with YouTube API",
	Long:  "Authenticate with YouTube API",
	Run: func(cmd *cobra.Command, args []string) {
		auth.NewY2BService(
			auth.WithCredential(credential),
			auth.WithCacheToken(cacheToken),
		)
	},
}

func init() {
	RootCmd.AddCommand(authCmd)

	authCmd.Flags().StringVarP(&credential, "credential", "c", "client_secret.json", "Path to client secret file")
	authCmd.Flags().StringVarP(&cacheToken, "cacheToken", "t", "youtube.token.json", "Path to token cache file")
}
