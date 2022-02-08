package certstore

import (
	"errors"
	"fmt"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certificate/service/factory"
	"bilalekrem.com/certstore/internal/certstore/config"
	"bilalekrem.com/certstore/internal/logging"
)

type certStoreImpl struct {
	certIssuers map[string]service.CertificateService
}

// -------

func NewFromConfig(conf *config.Config) (*certStoreImpl, error) {
	store := &certStoreImpl{
		certIssuers: make(map[string]service.CertificateService),
	}

	// ------

	for _, issuerConfig := range conf.IssuerConfigs {
		issuer := factory.NewService(issuerConfig.Type, issuerConfig.Args)

		store.RegisterIssuer(issuerConfig.Name, issuer)
	}

	return store, nil
}

// ------

func (c *certStoreImpl) IssueCertificate(issuer string, request *service.NewCertificateRequest) (*service.NewCertificateResponse, error) {
	certService, exist := c.certIssuers[issuer]
	if !exist {
		logging.GetLogger().Debug("Issuer not found: [%s]", issuer)
		return nil, errors.New(fmt.Sprintf("Issuer not found: [%s]", issuer))
	}

	// ----

	logging.GetLogger().Debugf("Issuer found, creating a new certificate %s", request)
	response, err := certService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ------

func (c *certStoreImpl) RegisterIssuer(issuer string, certService service.CertificateService) {
	logging.GetLogger().Debugf("Registering a new certificate service: [%s]", issuer)
	c.certIssuers[issuer] = certService
}
