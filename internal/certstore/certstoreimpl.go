package certstore

import (
	"errors"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/logging"
)

type certStoreImpl struct {
	caCertService service.CertificateService
	certService service.CertificateService
}

var DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS = 2 * 365

// -------

func New() (*certStoreImpl, error) {
	return &certStoreImpl{
		caCertService: &service.CACertificateService{},
	}, nil
}

func NewWithCA(caPrivateKeyPem []byte, caCertPem []byte) (*certStoreImpl, error) {
	certService, err := service.New(caPrivateKeyPem, caCertPem)
	if err != nil {
		return nil, err
	}

	return &certStoreImpl{
		caCertService: &service.CACertificateService{},
		certService: certService,
	}, nil
}

// ------

func (c *certStoreImpl) CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error) {
	request := &service.NewCertificateRequest{
		CommonName: clusterName,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
	}
	logging.GetLogger().Debug("creating cluster ca certificate")
	response, err := c.caCertService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *certStoreImpl) CreateServerCertificate(advertisedServerName string) (*service.NewCertificateResponse, error) {
	if c.certService == nil {
		return nil, errors.New("CA required to create and sign server ecertificates")
	}

	request := &service.NewCertificateRequest{
		CommonName: advertisedServerName,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
		SubjectAlternativeNames: []string{advertisedServerName},
	}

	logging.GetLogger().Debugf("creating server certificate for %s\n", advertisedServerName)
	response, err := c.certService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *certStoreImpl) CreateWorkerCertificate(address string) (*service.NewCertificateResponse, error) {
	if c.certService == nil {
		return nil, errors.New("CA required to create and sign worker ecertificates")
	}

	request := &service.NewCertificateRequest{
		CommonName: address,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
		SubjectAlternativeNames: []string{address},
	}

	logging.GetLogger().Debugf("creating worker certificate for %s\n", address)
	response, err := c.certService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}