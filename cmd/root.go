package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	short = "A fully functional CLI for YouTube"
	long  = "yutu is a fully functional CLI for YouTube, which can be used to manipulate YouTube videos, playlists, channels, etc"
)

var RootCmd = &cobra.Command{
	Use:   "yutu",
	Short: short,
	Long:  long,

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.AutomaticEnv()
}
