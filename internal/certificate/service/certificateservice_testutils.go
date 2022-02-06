package service

import (
	"regexp"
	"testing"
	"time"

	"crypto/x509"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/x509utils"
)

func testEmail(t *testing.T, service *CertificateService) {
	email := []string{"test@mail.com"}

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		Email:          email,
		ExpirationDays: 365,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	assert.NotEqualM(t, 0, cert.EmailAddresses, "email should've been found")
	assert.DeepEqual(t, email, cert.EmailAddresses)
}

func testNotValidEmail(t *testing.T, service *CertificateService) {
	email := []string{"test"}

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		Email:          email,
		ExpirationDays: 365,
	}
	_, err := (*service).CreateCertificate(request)

	assert.ErrorContains(t, err, "Validation error: email")
}

func testSubject(t *testing.T, service *CertificateService) {
	commonName := "my-ca"
	email := []string{"test@mail.com"}
	organization := []string{"my-org"}

	request := &NewCertificateRequest{
		CommonName:     commonName,
		Email:          email,
		Organization:   organization,
		ExpirationDays: 365,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	subject := cert.Subject
	assert.Equal(t, commonName, subject.CommonName)

	assert.NotEqual(t, 0, subject.Organization)
	assert.DeepEqual(t, organization, subject.Organization)
}

func testNotProvidedCommonName(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		ExpirationDays: 365,
	}
	_, err := (*service).CreateCertificate(request)

	assert.ErrorContains(t, err, "Validation error: common name")
}

func testUniqueSerialNumber(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	firstResponse := createCert(t, service, request)
	secondResponse := createCert(t, service, request)

	firstCert := parsePEMToX509Certificate(t, firstResponse.Certificate)
	secondCert := parsePEMToX509Certificate(t, secondResponse.Certificate)

	assert.NotEqualM(t, firstCert.SerialNumber, secondCert.SerialNumber,
		"Sequentially created two certs' serial numbers must be different")
}

func testExpirationDate(t *testing.T, service *CertificateService) {
	expirationDays := 742
	beforeExpirationDate := time.Now().AddDate(0, 0, expirationDays-5)

	// ----

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: expirationDays,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	// ---

	afterExpirationDate := time.Now().AddDate(0, 0, expirationDays+5)

	assert.True(t, cert.NotAfter.Before(afterExpirationDate))
	assert.True(t, beforeExpirationDate.Before(cert.NotAfter))
}

func testNotValidExpirationDate(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 0,
	}
	_, err := (*service).CreateCertificate(request)

	assert.ErrorContains(t, err, "Validation error: expiration days")
}

func testNotProvidedExpirationDate(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName: "my-ca",
	}

	_, err := (*service).CreateCertificate(request)
	assert.ErrorContains(t, err, "Validation error: expiration days")
}

func testStartDate(t *testing.T, service *CertificateService) {
	beforeCreateCert := time.Now().AddDate(0, 0, -1)

	// ----

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(t, response.Certificate)

	// ---

	afterCreateCert := time.Now()

	assert.True(t, beforeCreateCert.Before(cert.NotBefore))
	assert.True(t, afterCreateCert.After(cert.NotBefore))
	assert.True(t, afterCreateCert.Before(cert.NotAfter))
}

func testCertificate(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response := createCert(t, service, request)
	certificateStr := string(response.Certificate)

	pattern := "-----BEGIN CERTIFICATE-----(\n.*)*----END CERTIFICATE-----\n"
	re := regexp.MustCompile(pattern)
	matched := re.MatchString(certificateStr)
	assert.True(t, matched)
}

func testRSAPrivateKey(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response, _ := (*service).CreateCertificate(request)
	privateKeyStr := string(response.PrivateKey)

	pattern := "-----BEGIN RSA PRIVATE KEY-----(\n.*)*----END RSA PRIVATE KEY-----\n"
	re := regexp.MustCompile(pattern)
	matched := re.MatchString(privateKeyStr)
	assert.True(t, matched)
}

// ---

func parsePEMToX509Certificate(t *testing.T, certPem []byte) *x509.Certificate {
	cert, err := x509utils.ParsePemCertificate(certPem)
	assert.NotError(t, err, "parsing certificate failed")
	return cert
}

// -----

func createCert(t *testing.T, service *CertificateService, request *NewCertificateRequest) *NewCertificateResponse {
	response, err := (*service).CreateCertificate(request)
	assert.NotError(t, err, "cert creation failed")
	return response
}
