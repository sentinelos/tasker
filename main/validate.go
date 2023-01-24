package main

import (
	"errors"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/sentinelos/tasker/internal/constants"
	"github.com/sentinelos/tasker/internal/executor"
	"github.com/sentinelos/tasker/internal/taskfile"
	"github.com/spf13/cobra"
)

// NewValidateCmd represents the Task file validate.
func NewValidateCmd() *cobra.Command {
	var (
		flagFilename  string
		flagDirectory string
	)

	cmd := &cobra.Command{
		Use:     "validate",
		Aliases: []string{"v"},
		Short:   "Validate the Task file",
		Long:    `Validate the Task file`,
		Args:    cobra.ExactArgs(0),
		Example: `  ` + constants.AppName + ` validate -f main.tf
    validate the Task file from the specified filename

  ` + constants.AppName + ` validate -d .workflows
    validate the Task files from the specified directory`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(flagFilename) == 0 && len(flagDirectory) == 0 {
				return errors.New("please provide the Task file name or directory")
			}

			files := map[string]*hcl.File{}
			taskFile, diags := taskfile.LoadTaskFile(flagFilename, executor.NewContext())

			if taskFile != nil {
				files[flagFilename] = taskFile.Source
			}

			wr := hcl.NewDiagnosticTextWriter(os.Stdout, files, 80, true)
			return wr.WriteDiagnostics(diags)
		},
	}

	cmd.Flags().StringVarP(&flagFilename, "filename", "f", "", "name of the Task file")
	cmd.Flags().StringVarP(&flagDirectory, "directory", "d", "", "directory of the Task files")

	return cmd
}
