package service

import (
	"testing"

	"crypto/x509"

	"bilalekrem.com/certstore/internal/assert"
)

func TestCA_ISCA(t *testing.T) {
	var service CertificateService = createCACertificateService()
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	assert.True(t, cert.IsCA)
}

func TestCA_KeyUsageCertSign(t *testing.T) {
	var service CertificateService = createCACertificateService()
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	assert.EqualM(t, (cert.KeyUsage & x509.KeyUsageDigitalSignature), x509.KeyUsageDigitalSignature,
		"CA certificate does not have digital signature key usage")

	assert.EqualM(t, (cert.KeyUsage & x509.KeyUsageCertSign), x509.KeyUsageCertSign,
		"CA certificate does not have certificate sign key usage")
}

// ----- common certificate service tests

func TestCA_Email(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testEmail(t, &service)
}

func TestCA_NotValidEmail(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testNotValidEmail(t, &service)
}

func TestCA_Subject(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testSubject(t, &service)
}

func TestCA_NotProvidedCommonName(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testNotProvidedCommonName(t, &service)
}

func TestCA_UniqueSerialNumber(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testUniqueSerialNumber(t, &service)
}

func TestCA_ExpirationDate(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testExpirationDate(t, &service)
}

func TestCA_NotValidExpirationDate(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testNotValidExpirationDate(t, &service)
}

func TestCA_NotProvidedExpirationDate(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testNotProvidedExpirationDate(t, &service)
}

func TestCA_StartDate(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testStartDate(t, &service)
}

func TestCA_Certificate(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testCertificate(t, &service)
}

func TestCA_RSAPrivateKey(t *testing.T) {
	var service CertificateService = createCACertificateService()
	testRSAPrivateKey(t, &service)
}

// ------

func createCACertificateService() *CACertificateService {
	return &CACertificateService{}
}
