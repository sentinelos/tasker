package cmd

import (
	"github.com/spf13/cobra"

	"github.com/sentinelos/tasker/internal/version"
)

var shortVersion bool

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if shortVersion {
			version.PrintShortVersion()
		} else {
			version.PrintLongVersion()
		}
	},
}

func init() {
	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "Print the short version")
	rootCmd.AddCommand(versionCmd)
}
