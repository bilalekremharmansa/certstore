package certstore

import (
	"errors"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/logging"
)

type certStoreImpl struct {
	clusterService service.CertificateService
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
	}, nil
}

func NewWithoutCA() (*certStoreImpl, error) {
	return &certStoreImpl{}, nil
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