package certstore

import (
	"errors"
	"fmt"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/logging"
)

type certStoreImpl struct {
	clusterService service.CertificateService
	certIssuers map[string]service.CertificateService
}

var DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS = 2 * 365

// -------

func New(caPrivateKeyPem []byte, caCertPem []byte) (*certStoreImpl, error) {
	clusterService, err := service.New(caPrivateKeyPem, caCertPem)
	if err != nil {
		return nil, err
	}

	return &certStoreImpl{
		clusterService: clusterService,
		certIssuers: make(map[string]service.CertificateService),
	}, nil
}

func NewWithoutCA() (*certStoreImpl, error) {
	return &certStoreImpl{
		certIssuers: make(map[string]service.CertificateService),
	}, nil
}

// ------

func (*certStoreImpl) CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error) {
	request := &service.NewCertificateRequest{
		CommonName: clusterName,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
	}
	logging.GetLogger().Debug("creating cluster ca certificate")

	caCertificateService := &service.CACertificateService{}
	response, err := caCertificateService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *certStoreImpl) CreateServerCertificate(advertisedServerName string) (*service.NewCertificateResponse, error) {
	if c.clusterService == nil {
		return nil, errors.New("CA required to create and sign server ecertificates")
	}

	request := &service.NewCertificateRequest{
		CommonName: advertisedServerName,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
		SubjectAlternativeNames: []string{advertisedServerName},
	}

	logging.GetLogger().Debugf("creating server certificate for %s\n", advertisedServerName)
	response, err := c.clusterService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *certStoreImpl) CreateWorkerCertificate(address string) (*service.NewCertificateResponse, error) {
	if c.clusterService == nil {
		return nil, errors.New("CA required to create and sign worker ecertificates")
	}

	request := &service.NewCertificateRequest{
		CommonName: address,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
		SubjectAlternativeNames: []string{address},
	}

	logging.GetLogger().Debugf("creating worker certificate for %s\n", address)
	response, err := c.clusterService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *certStoreImpl) IssueCertificate(issuer string, request *service.NewCertificateRequest) (*service.NewCertificateResponse, error) {
	certService, exist := c.certIssuers[issuer]
	if !exist {
		logging.GetLogger().Debug("Issuer not found: [%s]", issuer)
		return nil, errors.New(fmt.Sprintf("Issuer not found: [%s]", issuer))
	}

	// ----

	logging.GetLogger().Debug("Issuer found, creating a new certificate %s", request)
	response, err := certService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ------

func (c *certStoreImpl) RegisterIssuer(issuer string, certService service.CertificateService) {
	logging.GetLogger().Debug("Registering a new certificate service: [%s]", issuer)
	c.certIssuers[issuer] = certService
}