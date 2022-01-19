package cli

import (
	wrk "bilalekrem.com/certstore/internal/cluster/worker"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

var clusterWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "manage worker",
}

var clusterWorkerCreateCertCmd = &cobra.Command{
	Use:   "createCert",
	Short: "create worker certificate and key",
	Run: func(cmd *cobra.Command, args []string) {
		advertisedName, _ := cmd.Flags().GetString("address")
		caCertPath, _ := cmd.Flags().GetString("cacert")
		caKeyPath, _ := cmd.Flags().GetString("cakey")

		// ---

		logging.GetLogger().Info("creating worker cert")
		certificate, err := getCertstoreWithCA(caKeyPath, caCertPath).CreateServerCertificate(advertisedName)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		saveCert(".", "worker", certificate)
	},
}

var clusterWorkerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start worker",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")
		pipelineToRun, _ := cmd.Flags().GetString("pipeline")

		// caCertPath, _ := cmd.Flags().GetString("cacert")
		// workerCertPath, _ := cmd.Flags().GetString("cert")
		// workerCertKeyPath, _ := cmd.Flags().GetString("certkey")

		// -----

		worker, err := wrk.NewFromFile(configPath)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		// ---

		err = worker.RunPipeline(pipelineToRun)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}
	},
}

func init() {
	clusterWorkerCreateCertCmd.Flags().String("address", "", "address of worker")
	clusterWorkerCreateCertCmd.Flags().String("cacert", "", "certificate authority file path in PEM format")
	clusterWorkerCreateCertCmd.Flags().String("cakey", "", "certificate authority key file path in PEM format")
	clusterWorkerCreateCertCmd.MarkFlagRequired("address")
	clusterWorkerCreateCertCmd.MarkFlagRequired("cacert")
	clusterWorkerCreateCertCmd.MarkFlagRequired("cakey")

	// ------
	clusterWorkerStartCmd.Flags().String("config", "", "worker config file path")
	clusterWorkerStartCmd.MarkFlagRequired("config")
	clusterWorkerStartCmd.Flags().String("pipeline", "", "pipeline name to run")
	clusterWorkerStartCmd.MarkFlagRequired("pipeline")

	// clusterWorkerStartCmd.Flags().String("cacert", "", "CA certificate file for verifying the server")
	// clusterWorkerStartCmd.Flags().String("cert", "", "x509 certificate file for mTLS")
	// clusterWorkerStartCmd.Flags().String("certkey", "", "x509 private key file for mTLS")
	// clusterWorkerStartCmd.MarkFlagRequired("cacert")
	// clusterWorkerStartCmd.MarkFlagRequired("cert")
	// clusterWorkerStartCmd.MarkFlagRequired("certkey")

	clusterCmd.AddCommand(clusterWorkerCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerCreateCertCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerStartCmd)
}
