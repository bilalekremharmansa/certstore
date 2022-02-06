package cli

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"path/filepath"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certstore"
	"bilalekrem.com/certstore/internal/certstore/config"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "manage cluster",
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}

// ----

func getCertstoreWithConfig(configPath string) certstore.CertStore {
	config, err := config.ParseFile(configPath)
	if err != nil {
		error("Error occurred while reading server config: [%v]\n", err)
	}

	store, err := certstore.NewFromConfig(config)
	if err != nil {
		error("creating certstore with config failed: [%v]\n", err)
	}

	return store
}

func saveCert(targetPath string, identifier string, certificate *service.NewCertificateResponse) {
	logging.GetLogger().Infof("Saving %s certificate and key", identifier)
	ioutil.WriteFile(filepath.Join(targetPath, identifier+".crt"), certificate.Certificate, 0644)
	ioutil.WriteFile(filepath.Join(targetPath, identifier+".key"), certificate.PrivateKey, 0600)
}

func createServerTLSConfig(caCertPath string, serverCertPath string, serverCertKeyPath string) *tls.Config {
	caCertPem, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		error("Error occurred while reading ca: [%v]\n", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCertPem) {
		error("could not add ca cert to cert pool")
	}

	serverCertificate, err := tls.LoadX509KeyPair(serverCertPath, serverCertKeyPath)
	if err != nil {
		error("Loading server certification failed: %v", err)
	}

	return &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{serverCertificate},
		ClientCAs:    caPool,
	}
}

func createWorkerTLSConfig(caCertPath string, workerCertPath string, workerCertKeyPath string) *tls.Config {
	caCertPem, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		error("Error occurred while reading ca: [%v]\n", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCertPem) {
		error("could not add ca cert to cert pool")
	}

	workerCertificate, err := tls.LoadX509KeyPair(workerCertPath, workerCertKeyPath)
	if err != nil {
		error("Loading server certification failed: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{workerCertificate},
		RootCAs:      caPool,
	}
}
