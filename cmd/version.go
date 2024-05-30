package cmd

import (
	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
)

var (
	Version    = ""
	Commit     = ""
	CommitDate = ""
	Os         = ""
	Arch       = ""
	repo       = "https://github.com/eat-pray-ai/yutu"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of yutu",
	Long:  "Show the version of yutu",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("yutu🐰 version %s", Version)
		if Os != "" && Arch != "" {
			cmd.Printf(" %s/%s", Os, Arch)
		}
		if Commit != "" && CommitDate != "" {
			cmd.Printf("\nbuild %s-%s", Commit, CommitDate)
		}

		cmd.Println()
		cmd.Println("Star🌟:", termlink.Link("Github/eat-pray-ai/yutu", repo))

	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
