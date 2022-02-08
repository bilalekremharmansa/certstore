package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	certstore_pkg "bilalekrem.com/certstore/internal/certstore"
	grpc_gen "bilalekrem.com/certstore/internal/certstore/grpc/gen"
	grpc_service "bilalekrem.com/certstore/internal/certstore/grpc/service"
	"bilalekrem.com/certstore/internal/cluster/server/config"
	"bilalekrem.com/certstore/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	certstore  certstore_pkg.CertStore
	grpcServer *grpc.Server
	listenPort int
}

func NewFromFile(path string) (*Server, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err := config.Parse(string(bytes))
	if err != nil {
		return nil, err
	}

	return NewFromConfig(config)
}

func NewFromConfig(conf *config.Config) (*Server, error) {
	err := validateConfig(conf)
	if err != nil {
		return nil, err
	}

	certstore, err := certstore_pkg.NewFromConfig(&conf.CertStore)
	if err != nil {
		return nil, err
	}

	grpcServer, err := createAndSetupGrpcServer(conf, certstore)
	if err != nil {
		return nil, err
	}

	server := &Server{
		certstore:  certstore,
		grpcServer: grpcServer,
		listenPort: conf.ListenPort,
	}

	return server, nil
}

func (s *Server) Serve() error {
	logging.GetLogger().Debugf("Starting to listening on localhost:%d", s.listenPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.listenPort))
	if err != nil {
		return fmt.Errorf("error occurred while listening port, %v", err)
	}

	return s.grpcServer.Serve(listen)
}

// -----

func validateConfig(conf *config.Config) error {
	if conf.TlsCACert == "" {
		return fmt.Errorf("tls-ca-cert is required argument")
	} else if conf.TlsServerCert == "" {
		return fmt.Errorf("tls-server-cert is required argument")
	} else if conf.TlsServerCertKey == "" {
		return fmt.Errorf("tls-server-cert-key is required argument")
	} else if conf.ListenPort == 0 {
		return fmt.Errorf("port is required argument, missing or provided zero")
	}

	// should we also validate cerstore config in here ?
	return nil
}

func createAndSetupGrpcServer(conf *config.Config, certstore certstore_pkg.CertStore) (*grpc.Server, error) {
	tlsConfig, err := createTlsConfig(conf)
	if err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(tlsConfig)
	if creds == nil {
		return nil, err
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
	grpc_gen.RegisterCertificateServiceServer(grpcServer, grpc_service.NewCertificateService(certstore))
	reflection.Register(grpcServer)

	return grpcServer, nil
}

func createTlsConfig(conf *config.Config) (*tls.Config, error) {
	caCertPem, err := ioutil.ReadFile(conf.TlsCACert)
	if err != nil {
		return nil, fmt.Errorf("Error occurred while reading ca: [%v]\n", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCertPem) {
		return nil, fmt.Errorf("could not add ca cert to cert pool")
	}

	serverCertificate, err := tls.LoadX509KeyPair(conf.TlsServerCert, conf.TlsServerCertKey)
	if err != nil {
		return nil, fmt.Errorf("Loading server certification failed: %v", err)
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{serverCertificate},
		ClientCAs:    caPool,
	}, nil
}
