package factory

import (
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certificate/service/letsencrypt"
	"bilalekrem.com/certstore/internal/logging"
)

type ServiceType string

const (
	Simple               ServiceType = "Simple"
	CertificateAuthority             = "CertificateAuthority"
	LetsEncrypt                      = "LetsEncrypt"
	Unknown                          = "Unknown"
)

func NewService(t ServiceType, args map[string]string) service.CertificateService {
	logging.GetLogger().Debugf("Creating new service with type [%s], with args: [%v]\n", t, args)

	switch t {
	case Simple:
		caPrivateKeyPath := args["private-key"]
		caCertificatePath := args["certificate"]

		caPrivateKey, err := ioutil.ReadFile(caPrivateKeyPath)
		if err != nil {
			logging.GetLogger().Errorf("reading private key failed, %v", err)
		}
		caCertificate, err := ioutil.ReadFile(caCertificatePath)
		if err != nil {
			logging.GetLogger().Errorf("reading certificate failed, %v", err)
		}

		svc, err := service.New([]byte(caPrivateKey), []byte(caCertificate))
		if err != nil {
			logging.GetLogger().Errorf("error occurred while creating new certificate service, %v", err)
			return nil
		}
		return svc
	case CertificateAuthority:
		svc := &service.CACertificateService{}
		return svc
	case LetsEncrypt:
		userEmail := args["email"]
		if userEmail == "" {
			logging.GetLogger().Errorf("email is required field for lets encrypt service")
			return nil
		}

		userPrivateKeyPath := args["private-key"]
		if userPrivateKeyPath == "" {
			logging.GetLogger().Errorf("private-key is required field for lets encrypt service")
			return nil
		}

		provider := args["provider"]
		if provider == "" {
			logging.GetLogger().Errorf("provider is required field for lets encrypt service")
			return nil
		}

		svc, err := letsencrypt.New(userEmail, userPrivateKeyPath, provider)
		if err != nil {
			logging.GetLogger().Errorf("error occurred while creating new lets encrypt certificate service, %v", err)
			return nil
		}

		return svc
	case Unknown:
	default:
		return nil
	}
	return nil
}
