package cli

import (
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

var clusterServerCmd = &cobra.Command{
	Use:   "server",
	Short: "manage server",
}

var clusterServerCreateCertCmd = &cobra.Command{
	Use:   "createCert",
	Short: "create server certificate and key",
	Run: func(cmd *cobra.Command, args []string) {
		advertisedName, _ := cmd.Flags().GetString("advertised-name")
		caCertPath, _ := cmd.Flags().GetString("cacert")
		caKeyPath, _ := cmd.Flags().GetString("cakey")

		// ---

		logging.GetLogger().Info("creating server cert")
		certificate, err := getCertstoreWithCA(caKeyPath, caCertPath).CreateServerCertificate(advertisedName)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		saveCert(".", "server", certificate)
	},
}

func init() {
	clusterServerCreateCertCmd.Flags().String("advertised-name", "", "advertised address of server")
	clusterServerCreateCertCmd.Flags().String("cacert", "", "certificate authority file path in PEM format")
	clusterServerCreateCertCmd.Flags().String("cakey", "", "certificate authority key file path in PEM format")
	clusterServerCreateCertCmd.MarkFlagRequired("advertised-name")
	clusterServerCreateCertCmd.MarkFlagRequired("cacert")
	clusterServerCreateCertCmd.MarkFlagRequired("cakey")

	clusterCmd.AddCommand(clusterServerCmd)
	clusterServerCmd.AddCommand(clusterServerCreateCertCmd)
}
