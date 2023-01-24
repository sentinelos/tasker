// Package main contains definitions of CLI commands.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sentinelos/tasker/internal/constants"
)

func main() {
	app := &cobra.Command{
		Use:   constants.AppName,
		Short: "Tasker is a task executor",
		Long: `Tasker is a task executor

All of tasker features can be driven through the various commands below.
For help with any of those, simply call them with --help.`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	app.AddCommand(
		NewVersionCmd(),
		NewValidateCmd(),
	)

	app.InitDefaultHelpCmd()

	if err := app.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
