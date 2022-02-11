package manager

import (
	"errors"

	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certificate/service/factory"
	"bilalekrem.com/certstore/internal/logging"
)

var DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS = 2 * 365

type clusterManagerImpl struct {
	clusterCertService service.CertificateService
}

// will be used for new cluster creation
func NewForInitialization() (*clusterManagerImpl, error) {
	return &clusterManagerImpl{}, nil
}

func NewFromFile(caCertPath string, caKeyPath string) (*clusterManagerImpl, error) {
	args := make(map[string]string)
	args["certificate"] = caCertPath
	args["private-key"] = caKeyPath
	certService := factory.NewService(factory.Simple, args)

	return &clusterManagerImpl{clusterCertService: certService}, nil
}

func NewFromCA(cert []byte, key []byte) (*clusterManagerImpl, error) {
	certService, err := service.New(key, cert)
	if err != nil {
		return nil, err
	}

	return &clusterManagerImpl{clusterCertService: certService}, nil
}

func (*clusterManagerImpl) CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error) {
	logging.GetLogger().Debug("creating cluster ca certificate")
	caCertificateService := factory.NewService(factory.CertificateAuthority, nil)
	request := &service.NewCertificateRequest{
		CommonName:     clusterName,
		ExpirationDays: DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
	}
	response, err := caCertificateService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *clusterManagerImpl) CreateServerCertificate(name string) (*service.NewCertificateResponse, error) {
	if c.clusterCertService == nil {
		return nil, errors.New("CA required to create and sign server certificates")
	}

	request := &service.NewCertificateRequest{
		CommonName:              name,
		ExpirationDays:          DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
		SubjectAlternativeNames: []string{name},
	}

	logging.GetLogger().Debugf("creating server certificate for %s", name)
	response, err := c.clusterCertService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *clusterManagerImpl) CreateAgentCertificate(name string) (*service.NewCertificateResponse, error) {
	if c.clusterCertService == nil {
		return nil, errors.New("CA required to create and sign agent ecertificates")
	}

	request := &service.NewCertificateRequest{
		CommonName:              name,
		ExpirationDays:          DEFAULT_CLUSTER_CERT_EXPIRATION_DAYS,
		SubjectAlternativeNames: []string{name},
	}

	logging.GetLogger().Debugf("creating agent certificate for %s", name)
	response, err := c.clusterCertService.CreateCertificate(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
