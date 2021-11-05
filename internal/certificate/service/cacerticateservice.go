package service

import (
	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/logging"

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
		logging.GetLogger().Debug("validating ca certificate request failed: [%v]", err)
		return nil, err
	}

	// ------

	serialNumber, err := x509utils.GetRandomCertificateSerialNumber()
	if err != nil {
		logging.GetLogger().Debug("creating cert serial number failed: [%v]", err)
		return nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   request.CommonName,
			Organization: request.Organization,
		},
		EmailAddresses:        request.Email,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, request.ExpirationDays),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// ----
	logging.GetLogger().Debug("Generating private key for CA")
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logging.GetLogger().Debug("generating ca private key failed: [%v]", err)
		return nil, err
	}

	logging.GetLogger().Debug("Creating CA key")
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		logging.GetLogger().Debug("creating ca cert failed: [%v]", err)
		return nil, err
	}

	// ----- pem encode

	logging.GetLogger().Debug("Encoding certificate and key")
	caPrivateKeyPem, caPem := x509utils.EncodePEMCertAndKey(caPrivateKey, caBytes)

	// ------

	response := &NewCertificateResponse{
		Certificate: caPem.Bytes(),
		PrivateKey:  caPrivateKeyPem.Bytes(),
	}

	return response, nil
}
