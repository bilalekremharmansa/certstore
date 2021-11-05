package service

import (
	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/logging"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"time"
)

type certificateServiceImpl struct {
	ca           *x509.Certificate
	caPrivateKey *rsa.PrivateKey
}

func New(privateKeyPem []byte, caPem []byte) (*certificateServiceImpl, error) {
	caCert, err := x509utils.ParsePemCertificate(caPem)
	if err != nil {
		return nil, err
	}

	caKey, err := x509utils.ParsePemPrivateKey(privateKeyPem)
	if err != nil {
		return nil, err
	}

	return &certificateServiceImpl{
		ca:           caCert,
		caPrivateKey: caKey,
	}, nil
}

func (service *certificateServiceImpl) CreateCertificate(request *NewCertificateRequest) (*NewCertificateResponse, error) {
	err := service.validate()
	if err != nil {
		logging.GetLogger().Debug("validating certificate service failed: [%v]", err)
		return nil, err
	}

	err = validateCertificateRequest(request)
	if err != nil {
		logging.GetLogger().Debug("validating certificate request failed: [%v]", err)
		return nil, err
	}

	// -----

	serialNumber, err := x509utils.GetRandomCertificateSerialNumber()
	if err != nil {
		logging.GetLogger().Debug("creating cert serial number failed: [%v]", err)
		return nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   request.CommonName,
			Organization: request.Organization,
		},
		EmailAddresses: request.Email,
		DNSNames:       request.SubjectAlternativeNames,
		NotBefore:      time.Now(),
		NotAfter:       time.Now().AddDate(0, 0, request.ExpirationDays),
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:       x509.KeyUsageDigitalSignature,
	}

	certPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logging.GetLogger().Debug("generating private key failed: [%v]", err)
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, service.ca, &certPrivateKey.PublicKey, service.caPrivateKey)
	if err != nil {
		logging.GetLogger().Debug("creating cert failed: [%v]", err)
		return nil, err
	}

	logging.GetLogger().Debug("Encoding certificate and key")
	certPrivateKeyPem, certPem := x509utils.EncodePEMCertAndKey(certPrivateKey, certBytes)
	response := &NewCertificateResponse{
		Certificate: certPem.Bytes(),
		PrivateKey:  certPrivateKeyPem.Bytes(),
	}

	return response, nil
}

func (service *certificateServiceImpl) validate() error {
	if service.ca == nil {
		return errors.New("Validation error: ca pem required to create certificates")
	}
	if service.caPrivateKey == nil {
		return errors.New("Validation error: ca private key pem required to create certificates")
	}

	return nil
}
