package service

import (
	"testing"

	"crypto/x509"

	"bilalekrem.com/certstore/internal/assert"
)

func TestDefault_NotISCA(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)
	assert.False(t, cert.IsCA)
}

func TestDefault_KeyUsage(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	assert.EqualM(t, (cert.KeyUsage & x509.KeyUsageDigitalSignature), x509.KeyUsageDigitalSignature,
		"CA certificate does not have digital signature key usage")
}

func TestDefault_SubjectAlternativeDNSNames(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	dnsNames := []string{"mysite.com", "mytest.com", "localhost"}
	request := &NewCertificateRequest{
		CommonName:              "my-ca",
		ExpirationDays:          5,
		SubjectAlternativeNames: dnsNames,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	assert.NotEqual(t, 0, len(cert.DNSNames))
	assert.DeepEqualM(t, dnsNames, dnsNames, "sans are not equal")
}

func TestDefault_VerifySignedWithCA(t *testing.T) {
	service := createCertificateServiceImpl(t)
	var polymorphicService CertificateService = service

	// -----

	dnsNames := []string{"mysite.com", "localhost"}
	request := &NewCertificateRequest{
		CommonName:              "my-ca",
		ExpirationDays:          5,
		SubjectAlternativeNames: dnsNames,
	}
	response := createCert(t, &polymorphicService, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	// ----

	roots := x509.NewCertPool()
	roots.AddCert(service.ca)
	opts := x509.VerifyOptions{
		Roots: roots,
	}

	_, err := cert.Verify(opts)
	assert.NotError(t, err, "verification of CA is failed\n")
}

// ----- common certificate service tests

func TestDefault_Email(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testEmail(t, &service)
}

func TestDefault_NotValidEmail(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testNotValidEmail(t, &service)
}

func TestDefault_Subject(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testSubject(t, &service)
}

func TestDefault_NotProvidedCommonName(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testNotProvidedCommonName(t, &service)
}

func TestDefault_UniqueSerialNumber(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testUniqueSerialNumber(t, &service)
}

func TestDefault_ExpirationDate(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testExpirationDate(t, &service)
}

func TestDefault_NotValidExpirationDate(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testNotValidExpirationDate(t, &service)
}

func TestDefault_NotProvidedExpirationDate(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testNotProvidedExpirationDate(t, &service)
}

func TestDefault_StartDate(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testStartDate(t, &service)
}

func TestDefault_Certificate(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testCertificate(t, &service)
}

func TestDefault_RSAPrivateKey(t *testing.T) {
	var service CertificateService = createCertificateServiceImpl(t)
	testRSAPrivateKey(t, &service)
}

// ------

func createCertificateServiceImpl(t *testing.T) *certificateServiceImpl {
	caCertService := CACertificateService{}
	caRequest := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	caResponse, _ := caCertService.CreateCertificate(caRequest)

	service, err := New(caResponse.PrivateKey, caResponse.Certificate)
	assert.NotError(t, err, "creating certificate service failed")
	return service
}
