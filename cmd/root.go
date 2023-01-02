// Package cmd contains definitions of CLI commands.
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sentinelos/tasker/pkg/constants"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   constants.AppName,
	Short: "Tasker is a task runner.",
	Long:  ``,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
