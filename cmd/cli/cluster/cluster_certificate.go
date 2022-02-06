package cluster

import (
	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/cluster/manager"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

func newCertificateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificate",
		Short: "creates server or worker certificate",
		Run: func(cmd *cobra.Command, args []string) {
			certType, _ := cmd.Flags().GetString("type")
			certName, _ := cmd.Flags().GetString("name")
			caCertPath, _ := cmd.Flags().GetString("cacert")
			caKeyPath, _ := cmd.Flags().GetString("cakey")

			if certType != "server" && certType != "worker" {
				cliutils.Error("cert type should be either server or worker")
			}

			// ---

			createAndSaveCert(certType, certName, caCertPath, caKeyPath)
		},
	}

	// ----

	cmd.Flags().String("type", "", "type of certificate, possible args: [server,worker]")
	cmd.Flags().String("name", "", "common name of certificate")
	cmd.Flags().String("cacert", "", "cluster certificate authority file path in PEM format")
	cmd.Flags().String("cakey", "", "cluster certificate authority key file path in PEM format")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("cacert")
	cmd.MarkFlagRequired("cakey")
	return cmd
}

func createAndSaveCert(certType string, name string, caCertPath string, caKeyPath string) {
	logging.GetLogger().Infof("creating certificate for %s : [%s]", certType, name)

	clusterManager, err := manager.NewFromFile(caCertPath, caKeyPath)
	cliutils.ValidateNotError(err)

	var certificate *service.NewCertificateResponse
	if certType == "server" {
		certificate, err = clusterManager.CreateServerCertificate(name)
		cliutils.ValidateNotError(err)
	} else {
		certificate, err = clusterManager.CreateWorkerCertificate(name)
		cliutils.ValidateNotError(err)
	}

	saveCert(".", certType, certificate)
}
