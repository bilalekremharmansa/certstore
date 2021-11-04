package service

import (
	"errors"
	"net/mail"

	"bilalekrem.com/certstore/internal/certificate/x509utils"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"
)

type CACertificateService struct {
}

func (service *CACertificateService) CreateCertificate(request *NewCertificateRequest) (*NewCertificateResponse, error) {
	err := validateCertificateRequest(request)
	if err != nil {
		return nil, err
	}

	// ------

	serialNumber, err := x509utils.GetRandomCertificateSerialNumber()
	if err != nil {
		return nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   request.CommonName,
			Organization: []string{request.Organization},
		},
		EmailAddresses:        []string{request.Email},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, request.ExpirationDays),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// ----

	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, err
	}

	// ----- pem encode

	caPrivateKeyPem, caPem := x509utils.EncodePEMCertAndKey(caPrivateKey, caBytes)

	// ------

	response := &NewCertificateResponse{
		Certificate: caPem.Bytes(),
		PrivateKey:  caPrivateKeyPem.Bytes(),
	}

	return response, nil
}

func validateCertificateRequest(req *NewCertificateRequest) error {
	if req.CommonName == "" {
		return errors.New("Validation error: common name can not be empty")
	}

	if req.Email != "" {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			return errors.New("Validation error: email is not valid")
		}
	}

	if req.ExpirationDays < 1 {
		return errors.New("Validation error: expiration days must be bigger than 1")
	}

	return nil
}
