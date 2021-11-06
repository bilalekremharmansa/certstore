package cli

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certstore"
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

func getCertstore() certstore.CertStore {
	certStore, err := certstore.New()
	if err != nil {
		error("Creation of certstore failed")
	}

	return certStore
}

func getCertstoreWithCA(caKeyPath string, caCertPath string) certstore.CertStore {
	caCertPem, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		error("Error occurred while reading ca cert: [%v]\n", err)
	}

	caKeyPem, err := ioutil.ReadFile(caKeyPath)
	if err != nil {
		error("Error occurred while reading ca cert: [%v]\n", err)
	}

	// ----

	certStore, err := certstore.NewWithCA(caKeyPem, caCertPem)
	if err != nil {
		error("Creation of certstore with CA failed, [%v]\n", err)
	}

	return certStore
}

func saveCert(targetPath string, identifier string, certificate *service.NewCertificateResponse) {
	logging.GetLogger().Infof("Saving %s certificate and key", identifier)
	ioutil.WriteFile(targetPath+"/"+identifier+".crt", certificate.Certificate, 0644)
	ioutil.WriteFile(targetPath+"/"+identifier+".key", certificate.PrivateKey, 0600)
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