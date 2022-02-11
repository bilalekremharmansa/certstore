package manager

import "bilalekrem.com/certstore/internal/certificate/service"

type ClusterManager interface {
	CreateClusterCACertificate(clusterName string) (*service.NewCertificateResponse, error)

	CreateServerCertificate(advertisedServerName string) (*service.NewCertificateResponse, error)
	CreateAgentCertificate(address string) (*service.NewCertificateResponse, error)
}
