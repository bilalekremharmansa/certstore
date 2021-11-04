package service

import (
	"testing"

	"crypto/x509"
)

func TestCA_ISCA(t *testing.T) {
	var service CertificateService = createCertificateService()
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(response.Certificate)
	if !cert.IsCA {
		t.Fatal("Certificate is not CA")
	}
}

func TestCA_KeyUsageCertSign(t *testing.T) {
	var service CertificateService = createCertificateService()
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, &service, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	if (cert.KeyUsage & x509.KeyUsageDigitalSignature) != x509.KeyUsageDigitalSignature {
		t.Fatalf("CA certificate does not have digital signature key usage, should've been")
	}

	if (cert.KeyUsage & x509.KeyUsageCertSign) != x509.KeyUsageCertSign {
		t.Fatalf("CA certificate does not have certificate sign key usage, should've been")
	}
}

// ----- common certificate service tests

func TestCA_Email(t *testing.T) {
	var service CertificateService = createCertificateService()
	testEmail(t, &service)
}

func TestCA_NotValidEmail(t *testing.T) {
	var service CertificateService = createCertificateService()
	testNotValidEmail(t, &service)
}

func TestCA_Subject(t *testing.T) {
	var service CertificateService = createCertificateService()
	testSubject(t, &service)
}

func TestCA_NotProvidedCommonName(t *testing.T) {
	var service CertificateService = createCertificateService()
	testNotProvidedCommonName(t, &service)
}

func TestCA_UniqueSerialNumber(t *testing.T) {
	var service CertificateService = createCertificateService()
	testUniqueSerialNumber(t, &service)
}

func TestCA_ExpirationDate(t *testing.T) {
	var service CertificateService = createCertificateService()
	testExpirationDate(t, &service)
}

func TestCA_NotValidExpirationDate(t *testing.T) {
	var service CertificateService = createCertificateService()
	testNotValidExpirationDate(t, &service)
}

func TestCA_NotProvidedExpirationDate(t *testing.T) {
	var service CertificateService = createCertificateService()
	testNotProvidedExpirationDate(t, &service)
}

func TestCA_StartDate(t *testing.T) {
	var service CertificateService = createCertificateService()
	testStartDate(t, &service)
}

func TestCA_Certificate(t *testing.T) {
	var service CertificateService = createCertificateService()
	testCertificate(t, &service)
}

func TestCA_RSAPrivateKey(t *testing.T) {
	var service CertificateService = createCertificateService()
	testRSAPrivateKey(t, &service)
}

// ------

func createCertificateService() *CACertificateService {
	return &CACertificateService{}
}
