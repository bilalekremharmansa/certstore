package cli

import (
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

var clusterInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clusterName, _ := cmd.Flags().GetString("name")

		// ---

		logging.GetLogger().Infof("creating cluster CA, cluster name: [%s]", clusterName)
		certificate, err := getCertstore().CreateClusterCACertificate(clusterName)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		saveCert(".", "ca", certificate)
	},
}

func init() {
	clusterInitCmd.Flags().String("name", "", "cluster name")
	clusterInitCmd.MarkFlagRequired("name")

	clusterCmd.AddCommand(clusterInitCmd)
}
