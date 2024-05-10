package cmd

import "github.com/spf13/cobra"

var (
	Version    = ""
	Commit     = ""
	CommitDate = ""
	Os         = ""
	Arch       = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of yutu",
	Long:  "Show the version of yutu",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("yutu version %s %s/%s\nbuild %s-%s\n", Version, Os, Arch, Commit, CommitDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
