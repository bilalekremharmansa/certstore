package letsencrypt

import (
	"errors"
	"os"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/lego"
	"bilalekrem.com/certstore/internal/lego/provider"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge"
	real_lego "github.com/go-acme/lego/lego"
)

type letsEncryptCertificateService struct {
	lego lego.LegoAdapter
}

func New(email string, privateKeyPath string, providerName string) (*letsEncryptCertificateService, error) {
	provider, err := getProvider(providerName)
	if err != nil {
		return nil, err
	}

	var adapter lego.LegoAdapter

	// ---

	_, err = os.OpenFile(privateKeyPath, os.O_RDONLY, 0666)
	if errors.Is(err, os.ErrNotExist) {
		logging.GetLogger().Warn("acme user private key path is not found, generating a new user")
		adapter, err = lego.NewAdapterWithNewUserRegistration(email, privateKeyPath, provider, real_lego.LEDirectoryStaging)
		if err != nil {
			return nil, err
		}
	} else {
		user, err := lego.NewAcmeUserWithPrivateKeyFile(email, privateKeyPath)
		if err != nil {
			return nil, err
		}

		adapter, err = lego.NewAdapter(user, provider, real_lego.LEDirectoryStaging)
		if err != nil {
			return nil, err
		}
	}

	return &letsEncryptCertificateService{lego: adapter}, nil
}

func (c *letsEncryptCertificateService) CreateCertificate(request *service.NewCertificateRequest) (*service.NewCertificateResponse, error) {
	logging.GetLogger().Info("Creating certificate with lets encrypt service")
	logging.GetLogger().Warnf("Lets encrpyt certificate service ignores 'email', 'organization', 'expiration days' fields")

	// ----

	err := validateCertificateRequest(request)
	if err != nil {
		logging.GetLogger().Debug("validating certificate request failed: [%v]", err)
		return nil, err
	}

	// ----

	domains := []string{request.CommonName}
	for _, san := range request.SubjectAlternativeNames {
		domains = append(domains, san)
	}

	obtainRequest := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	obtainResource, err := c.lego.Obtain(obtainRequest)
	if err != nil {
		return nil, err
	}
	cert := obtainResource.Certificate
	privateKey := obtainResource.PrivateKey

	return &service.NewCertificateResponse{
		Certificate: cert,
		PrivateKey:  privateKey,
	}, nil
}

// ---

func validateCertificateRequest(req *service.NewCertificateRequest) error {
	if req.CommonName == "" {
		return errors.New("Validation error: common name can not be empty")
	}

	return nil
}

func getProvider(providerName string) (challenge.Provider, error) {
	if providerName == "mock" {
		return &provider.MockDNSProvider{}, nil
	}

	logging.GetLogger().Errorf("lego challenge not found with name, %s", providerName)
	return nil, errors.New("lego challenge not found with name")
}
