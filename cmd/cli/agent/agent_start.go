package worker

import (
	"runtime"

	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	wrk "bilalekrem.com/certstore/internal/cluster/worker"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start worker",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")

			// -----

			_, err := wrk.NewFromFile(configPath)
			cliutils.ValidateNotError(err)

			// ---

			runtime.Goexit()
		},
	}

	cmd.Flags().String("config", "", "worker config file path")
	cmd.MarkFlagRequired("config")
	return cmd
}
