package cluster

import (
	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	"bilalekrem.com/certstore/internal/cluster/manager"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init cluster",
		Run: func(cmd *cobra.Command, args []string) {
			clusterName, _ := cmd.Flags().GetString("name")

			// ---

			clusterManager, err := manager.NewForInitialization()
			cliutils.ValidateNotError(err)

			logging.GetLogger().Infof("creating cluster CA, cluster name: [%s]", clusterName)
			certificate, err := clusterManager.CreateClusterCACertificate(clusterName)
			cliutils.ValidateNotError(err)

			saveCert(".", "ca", certificate)
		},
	}

	// ----

	cmd.Flags().String("name", "", "cluster name")
	cmd.MarkFlagRequired("name")
	return cmd
}
