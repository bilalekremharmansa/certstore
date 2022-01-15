package factory

import (
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/logging"
)

type ServiceType int

const (
	Simple ServiceType = 1 << iota
	CertificateAuthority
)

func NewService(t ServiceType, args map[string]interface{}) service.CertificateService {
	switch t {
	case Simple:
		caPrivateKeyPath := args["private-key"].(string)
		caCertificatePath := args["certificate"].(string)

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
	default:
		return nil
	}
}
