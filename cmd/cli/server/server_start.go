package server

import (
	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	cluster_server_pkg "bilalekrem.com/certstore/internal/cluster/server"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.Flags().GetString("config")

			// ----

			server, err := cluster_server_pkg.NewFromFile(configPath)
			cliutils.ValidateNotError(err)

			server.Serve()
		},
	}

	cmd.Flags().String("config", "", "agent config file path")
	cmd.MarkFlagRequired("config")
	return cmd
}
