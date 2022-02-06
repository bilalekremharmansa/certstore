package cli

import (
	"runtime"

	"bilalekrem.com/certstore/internal/cluster/manager"
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
		name, _ := cmd.Flags().GetString("name")
		caCertPath, _ := cmd.Flags().GetString("cacert")
		caKeyPath, _ := cmd.Flags().GetString("cakey")

		// ---

		clusterManager, err := manager.NewFromFile(caKeyPath, caCertPath)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		logging.GetLogger().Info("creating worker cert")
		certificate, err := clusterManager.CreateWorkerCertificate(name)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		saveCert(".", "worker", certificate)
	},
}

var clusterWorkerRunPipelineCmd = &cobra.Command{
	Use:   "runPipeline",
	Short: "run a pipeline",
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

var clusterWorkerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start worker",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := cmd.Flags().GetString("config")

		// caCertPath, _ := cmd.Flags().GetString("cacert")
		// workerCertPath, _ := cmd.Flags().GetString("cert")
		// workerCertKeyPath, _ := cmd.Flags().GetString("certkey")

		// -----

		_, err := wrk.NewFromFile(configPath)
		if err != nil {
			error("error occurred: [%v]\n", err)
		}

		// ---

		runtime.Goexit()
	},
}

func init() {
	clusterWorkerCreateCertCmd.Flags().String("name", "", "name of worker")
	clusterWorkerCreateCertCmd.Flags().String("cacert", "", "certificate authority file path in PEM format")
	clusterWorkerCreateCertCmd.Flags().String("cakey", "", "certificate authority key file path in PEM format")
	clusterWorkerCreateCertCmd.MarkFlagRequired("name")
	clusterWorkerCreateCertCmd.MarkFlagRequired("cacert")
	clusterWorkerCreateCertCmd.MarkFlagRequired("cakey")

	// ------
	clusterWorkerRunPipelineCmd.Flags().String("config", "", "worker config file path")
	clusterWorkerRunPipelineCmd.MarkFlagRequired("config")
	clusterWorkerRunPipelineCmd.Flags().String("pipeline", "", "pipeline name to run")
	clusterWorkerRunPipelineCmd.MarkFlagRequired("pipeline")

	// clusterWorkerRunPipelineCmd.Flags().String("cacert", "", "CA certificate file for verifying the server")
	// clusterWorkerRunPipelineCmd.Flags().String("cert", "", "x509 certificate file for mTLS")
	// clusterWorkerRunPipelineCmd.Flags().String("certkey", "", "x509 private key file for mTLS")
	// clusterWorkerRunPipelineCmd.MarkFlagRequired("cacert")
	// clusterWorkerRunPipelineCmd.MarkFlagRequired("cert")
	// clusterWorkerRunPipelineCmd.MarkFlagRequired("certkey")

	// ------

	clusterWorkerStartCmd.Flags().String("config", "", "worker config file path")
	clusterWorkerStartCmd.MarkFlagRequired("config")

	clusterCmd.AddCommand(clusterWorkerCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerCreateCertCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerRunPipelineCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerStartCmd)
}
