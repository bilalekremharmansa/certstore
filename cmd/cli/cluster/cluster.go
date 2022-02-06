package cluster

import (
	"io/ioutil"
	"path/filepath"

	cliutils "bilalekrem.com/certstore/cmd/cli/utils"
	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "manage cluster",
	}

	// ----

	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newCertificateCommand())
	return cmd
}

// ----

func saveCert(targetPath string, identifier string, certificate *service.NewCertificateResponse) {
	certPath := filepath.Join(targetPath, identifier+".crt")
	logging.GetLogger().Infof("Saving %s certificate: [%s]", identifier, certPath)
	err := ioutil.WriteFile(certPath, certificate.Certificate, 0644)
	cliutils.ValidateNotError(err)

	keyPath := filepath.Join(targetPath, identifier+".key")
	logging.GetLogger().Infof("Saving %s private key: [%s]", identifier, keyPath)
	ioutil.WriteFile(keyPath, certificate.PrivateKey, 0600)
	cliutils.ValidateNotError(err)
}
