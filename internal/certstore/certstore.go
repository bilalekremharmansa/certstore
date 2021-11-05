package certstore

import "bilalekrem.com/certstore/internal/certificate/service"

type CertStore interface {
	CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error)
}