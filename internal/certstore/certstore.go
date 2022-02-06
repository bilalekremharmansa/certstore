package certstore

import "bilalekrem.com/certstore/internal/certificate/service"

type CertStore interface {
	IssueCertificate(string, *service.NewCertificateRequest) (*service.NewCertificateResponse, error)
}
