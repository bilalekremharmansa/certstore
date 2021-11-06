package certstore

import "bilalekrem.com/certstore/internal/certificate/service"

type CertStore interface {
	CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error)

	CreateServerCertificate(advertisedServerName string) (*service.NewCertificateResponse, error)
	CreateWorkerCertificate(address string) (*service.NewCertificateResponse, error)
}