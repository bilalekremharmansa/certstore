package server

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "manage server",
	}

	// ----

	cmd.AddCommand(newStartCommand())
	return cmd
}
