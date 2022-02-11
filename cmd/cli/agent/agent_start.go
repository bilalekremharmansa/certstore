package agent

import (
	"runtime"

	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	wrk "bilalekrem.com/certstore/internal/cluster/agent"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start agent",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")

			// -----

			_, err := wrk.NewFromFile(configPath)
			cliutils.ValidateNotError(err)

			// ---

			runtime.Goexit()
		},
	}

	cmd.Flags().String("config", "", "agent config file path")
	cmd.MarkFlagRequired("config")
	return cmd
}
