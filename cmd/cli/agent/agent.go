package agent

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "manage agent",
	}

	// ----

	cmd.AddCommand(newRunPipelineCommand())
	cmd.AddCommand(newStartCommand())
	return cmd
}
