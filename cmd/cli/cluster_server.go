package cli

import (
	"fmt"
	"net"

	pb "bilalekrem.com/certstore/internal/grpc/proto"
	"bilalekrem.com/certstore/internal/grpc/server"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
		caCertPath, _ := cmd.Flags().GetString("cacert")
		serverCertPath, _ := cmd.Flags().GetString("cert")
		serverCertKeyPath, _ := cmd.Flags().GetString("certkey")

		// ----

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			error("error occurred while listening port, %v", err)
		}

		tlsConfig := createServerTLSConfig(caCertPath, serverCertPath, serverCertKeyPath)
		creds := credentials.NewTLS(tlsConfig)
		if creds == nil {
			error("Failed to generate credentials %v", err)
		}

		opts := []grpc.ServerOption{grpc.Creds(creds)}
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
	clusterServerStartCmd.Flags().String("cacert", "", "CA certificate file for verifying the server")
	clusterServerStartCmd.Flags().String("cert", "", "x509 certificate file for mTLS")
	clusterServerStartCmd.Flags().String("certkey", "", "x509 private key file for mTLS")
	clusterServerStartCmd.MarkFlagRequired("cacert")
	clusterServerStartCmd.MarkFlagRequired("cert")
	clusterServerStartCmd.MarkFlagRequired("certkey")

	// ------

	clusterCmd.AddCommand(clusterServerCmd)
	clusterServerCmd.AddCommand(clusterServerCreateCertCmd)
	clusterServerCmd.AddCommand(clusterServerStartCmd)
}
