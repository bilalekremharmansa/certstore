package service

type NewCertificateRequest struct {
	CommonName     string
	Email          string
	Organization   string
	ExpirationDays int

	SubjectAlternativeNames []string
}

type NewCertificateResponse struct {
	// both, certificate and private key, is encoded in PEM format.
	Certificate []byte
	PrivateKey  []byte
}

type CertificateService interface {
	CreateCertificate(*NewCertificateRequest) (*NewCertificateResponse, error)
}
