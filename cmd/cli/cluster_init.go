package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

var clusterInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init cluster",
	Run: func(cmd *cobra.Command, args []string) {
		logging.GetLogger().Info("creating cluster CA")
		certificate, err := certStore.CreateClusterCACertificate(clusterName)
		if err != nil {
			fmt.Printf("error occurred: [%v]", err)
			os.Exit(1)
		}

		logging.GetLogger().Info("Saving CA and key")
		ioutil.WriteFile("ca.crt", certificate.Certificate, 0644)
		ioutil.WriteFile("ca.key", certificate.PrivateKey, 0644)
	},
}

func init() {
	clusterInitCmd.Flags().StringVar(&clusterName, "name", "", "cluster name")
	clusterInitCmd.MarkFlagRequired("name")

	clusterCmd.AddCommand(clusterInitCmd)
}
