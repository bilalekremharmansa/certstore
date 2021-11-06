package cli

import (
	"fmt"
	"net"

	pb "bilalekrem.com/certstore/internal/grpc/proto"
	"bilalekrem.com/certstore/internal/grpc/server"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

var clusterServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		// ----

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			error("error occurred while listening port, %v", err)
		}
		var opts []grpc.ServerOption
		grpcServer := grpc.NewServer(opts...)
		pb.RegisterHelloServiceServer(grpcServer, server.NewHelloService())
		reflection.Register(grpcServer)
		grpcServer.Serve(lis)
	},
}

func init() {
	clusterServerCreateCertCmd.Flags().String("advertised-name", "", "advertised address of server")
	clusterServerCreateCertCmd.Flags().String("cacert", "", "certificate authority file path in PEM format")
	clusterServerCreateCertCmd.Flags().String("cakey", "", "certificate authority key file path in PEM format")
	clusterServerCreateCertCmd.MarkFlagRequired("advertised-name")
	clusterServerCreateCertCmd.MarkFlagRequired("cacert")
	clusterServerCreateCertCmd.MarkFlagRequired("cakey")

	// ------

	clusterServerStartCmd.Flags().Int("port", 10000, "listen port")

	// ------

	clusterCmd.AddCommand(clusterServerCmd)
	clusterServerCmd.AddCommand(clusterServerCreateCertCmd)
	clusterServerCmd.AddCommand(clusterServerStartCmd)
}
