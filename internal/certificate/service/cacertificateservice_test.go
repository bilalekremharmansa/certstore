package service

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"crypto/x509"
	"encoding/pem"
)

func TestISCA(t *testing.T) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response := createCert(t, request)
	cert := parsePEMToX509Certificate(response.Certificate)
	if !cert.IsCA {
		t.Fatal("Certificate is not CA")
	}
}

func TestEmail(t *testing.T) {
	email := "test@mail.com"

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		Email:          email,
		ExpirationDays: 365,
	}
	response := createCert(t, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	if len(cert.EmailAddresses) == 0 || cert.EmailAddresses[0] != email {
		t.Fatalf("Email address is not correct, expected: [%s] found: [%s] \n", email, cert.EmailAddresses)
	}
}

func TestNotValidEmail(t *testing.T) {
	email := "test"

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		Email:          email,
		ExpirationDays: 365,
	}
	caCertService := CACertificateService{}
	_, err := caCertService.CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: email")) {
		t.Fatalf("provided not valid email, but expected error not raised error: [%v]", err)
	}
}

func TestSubject(t *testing.T) {
	commonName := "my-ca"
	email := "test@mail.com"
	organization := "my-org"

	request := &NewCertificateRequest{
		CommonName:     commonName,
		Email:          email,
		Organization:   organization,
		ExpirationDays: 365,
	}
	response := createCert(t, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	subject := cert.Subject
	if subject.CommonName != commonName {
		t.Fatalf("Common name is not correct, expected: [%s] found: [%s] \n", commonName, subject.CommonName)
	}

	if len(subject.Organization) == 0 || subject.Organization[0] != organization {
		t.Fatalf("Organization name is not correct, expected: [%s] found: [%s] \n", organization, subject.Organization)
	}
}

func TestNotProvidedCommonName(t *testing.T) {
	request := &NewCertificateRequest{
		ExpirationDays: 365,
	}
	caCertService := CACertificateService{}
	_, err := caCertService.CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: common name")) {
		t.Fatalf("provided not valid email, but error proper error not raised error: [%v]", err)
	}
}

func TestUniqueSerialNumber(t *testing.T) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	firstResponse := createCert(t, request)
	secondResponse := createCert(t, request)

	firstCert := parsePEMToX509Certificate(firstResponse.Certificate)
	secondCert := parsePEMToX509Certificate(secondResponse.Certificate)

	if firstCert.SerialNumber == secondCert.SerialNumber {
		t.Fatalf("Sequentially created two certs' serial numbers are same, should've been different. first: [%s], second: [%s]",
			firstCert.SerialNumber.String(), secondCert.SerialNumber.String())
	}
}

func TestExpirationDate(t *testing.T) {
	expirationDays := 742
	beforeExpirationDate := time.Now().AddDate(0, 0, expirationDays)

	// ----

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: expirationDays,
	}
	response := createCert(t, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	// ---

	certExpireDate := cert.NotAfter
	afterExpirationDate := time.Now().AddDate(0, 0, expirationDays)

	if certExpireDate.After(beforeExpirationDate) && certExpireDate.Before(afterExpirationDate) {
		t.Fatalf("Expiration date is not correct, expected [%d] days later, found: [%s] \n", expirationDays, certExpireDate)
	}
}

func TestNotValidExpirationDate(t *testing.T) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 0,
	}
	caCertService := CACertificateService{}
	_, err := caCertService.CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: expiration days")) {
		t.Fatalf("provided not valid expiration date, must be bigger than 1: [%v]", err)
	}
}

func TestNotProvidedExpirationDate(t *testing.T) {
	request := &NewCertificateRequest{
		CommonName: "my-ca",
	}
	caCertService := CACertificateService{}
	_, err := caCertService.CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: expiration days")) {
		t.Fatalf("provided not valid expiration date, must be bigger than 1: [%v]", err)
	}
}

func TestStartDate(t *testing.T) {
	beforeCreateCert := time.Now()

	// ----

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	// ---

	certExpireDate := cert.NotBefore
	afterCreateCert := time.Now()

	if certExpireDate.After(beforeCreateCert) && cert.NotAfter.Before(afterCreateCert) {
		t.Fatalf("Start date is not correct, should've been now but: [%s] \n", certExpireDate)
	}
}

func TestKeyUsageCertSign(t *testing.T) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	if (cert.KeyUsage & x509.KeyUsageDigitalSignature) != x509.KeyUsageDigitalSignature {
		t.Fatalf("CA certificate does not have digital signature key usage, should've been")
	}

	if (cert.KeyUsage & x509.KeyUsageCertSign) != x509.KeyUsageCertSign {
		t.Fatalf("CA certificate does not have certificate sign key usage, should've been")
	}
}

func TestCertificate(t *testing.T) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response := createCert(t, request)
	certificateStr := string(response.Certificate)

	pattern := "-----BEGIN CERTIFICATE-----(\n.*)*----END CERTIFICATE-----\n"
	re := regexp.MustCompile(pattern)
	matched := re.MatchString(certificateStr)
	if !matched {
		t.Fatalf("expected [%s] but found [%s]", pattern, certificateStr)
	}
}

func createCert(t *testing.T, request *NewCertificateRequest) *NewCertificateResponse {
	caCertService := CACertificateService{}
	response, err := caCertService.CreateCertificate(request)
	if err != nil {
		t.Fatalf("cert creation failed %v", err)
	}
	return response
}

func TestPrivateKey(t *testing.T) {
	caCertService := CACertificateService{}

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	response, _ := caCertService.CreateCertificate(request)
	privateKeyStr := string(response.PrivateKey)

	pattern := "-----BEGIN RSA PRIVATE KEY-----(\n.*)*----END RSA PRIVATE KEY-----\n"
	re := regexp.MustCompile(pattern)
	matched := re.MatchString(privateKeyStr)
	if !matched {
		t.Fatalf("expected [%s] but found [%s]", pattern, privateKeyStr)
	}
}

// -----

func parsePEMToX509Certificate(certPem []byte) *x509.Certificate {
	block, _ := pem.Decode(certPem)
	if block == nil {
		panic("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	return cert
}
