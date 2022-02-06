package server

import (
	"fmt"
	"net"

	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	"bilalekrem.com/certstore/internal/certstore"
	"bilalekrem.com/certstore/internal/certstore/config"
	grpc_gen "bilalekrem.com/certstore/internal/certstore/grpc/gen"
	grpc_service "bilalekrem.com/certstore/internal/certstore/grpc/service"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			certStoreConfigPath, _ := cmd.Flags().GetString("config")

			// ----

			certstore := getCertstoreWithConfig(certStoreConfigPath)

			// ----

			logging.GetLogger().Debugf("Starting to listening on localhost:%d", port)
			lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
			cliutils.ValidateNotError(err)

			opts := []grpc.ServerOption{}
			grpcServer := grpc.NewServer(opts...)
			grpc_gen.RegisterCertificateServiceServer(grpcServer, grpc_service.NewCertificateService(certstore))
			reflection.Register(grpcServer)
			grpcServer.Serve(lis)
		},
	}

	cmd.Flags().Int("port", 10000, "listen port")

	cmd.Flags().String("config", "", "worker config file path")
	cmd.MarkFlagRequired("config")
	return cmd
}

func getCertstoreWithConfig(configPath string) certstore.CertStore {
	config, err := config.ParseFile(configPath)
	cliutils.ValidateNotError(err)

	store, err := certstore.NewFromConfig(config)
	cliutils.ValidateNotError(err)
	return store
}
