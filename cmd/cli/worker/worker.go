package worker

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "manage worker",
	}

	// ----

	cmd.AddCommand(newRunPipelineCommand())
	cmd.AddCommand(newStartCommand())
	return cmd
}
