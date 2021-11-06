package cli

import (
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
