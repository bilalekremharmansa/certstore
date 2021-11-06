package cli

import (
	"context"
	"fmt"

	pb "bilalekrem.com/certstore/internal/grpc/proto"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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
	Use:   "communicateServer",
	Short: "communicat with server",
	Run: func(cmd *cobra.Command, args []string) {
		serverAddress, _ := cmd.Flags().GetString("server-address")

		// -----

		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())

		conn, err := grpc.Dial(serverAddress, opts...)
		if err != nil {
			error("fail to dial gRPC server: %v", err)
		}
		defer conn.Close()

		client := pb.NewHelloServiceClient(conn)

		request := &pb.HelloRequest{
			Name: "certstore!",
		}
		response, _ := client.SayHello(context.Background(), request)
		fmt.Println(response.Message)
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
	clusterWorkerStartCmd.Flags().String("server-address", "", "server address to communicate")
	clusterWorkerStartCmd.MarkFlagRequired("server-address")

	clusterCmd.AddCommand(clusterWorkerCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerCreateCertCmd)
	clusterWorkerCmd.AddCommand(clusterWorkerStartCmd)
}
