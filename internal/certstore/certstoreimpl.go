package certstore

import "bilalekrem.com/certstore/internal/certificate/service"

type certStoreImpl struct {
	caCertService *service.CACertificateService
}

var DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS = 720

// -------

func New() (*certStoreImpl, error) {
	return &certStoreImpl{
		caCertService: &service.CACertificateService{},
	}, nil
}

func (c *certStoreImpl) CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error) {
	request := &service.NewCertificateRequest{
		CommonName: clusterName,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
	}
	response, err := c.caCertService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}