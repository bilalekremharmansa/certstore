package factory

import (
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/logging"
)

type ServiceType string

const (
	Simple               ServiceType = "Simple"
	CertificateAuthority             = "CertificateAuthority"
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
	case Unknown:
	default:
		return nil
	}
	return nil
}
