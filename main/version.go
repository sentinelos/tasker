package main

import (
	"github.com/spf13/cobra"

	"github.com/sentinelos/tasker/internal/version"
)

// NewVersionCmd represents the version command.
func NewVersionCmd() *cobra.Command {
	var shortVersion bool

	cmd := &cobra.Command{
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

	cmd.Flags().BoolVar(&shortVersion, "short", false, "Print the short version")

	return cmd
}
