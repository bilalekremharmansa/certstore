package service

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"crypto/x509"
	"encoding/pem"
)

func testEmail(t *testing.T, service *CertificateService) {
	email := "test@mail.com"

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		Email:          email,
		ExpirationDays: 365,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	if len(cert.EmailAddresses) == 0 || cert.EmailAddresses[0] != email {
		t.Fatalf("Email address is not correct, expected: [%s] found: [%s] \n", email, cert.EmailAddresses)
	}
}

func testNotValidEmail(t *testing.T, service *CertificateService) {
	email := "test"

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		Email:          email,
		ExpirationDays: 365,
	}
	_, err := (*service).CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: email")) {
		t.Fatalf("provided not valid email, but expected error not raised error: [%v]", err)
	}
}

func testSubject(t *testing.T, service *CertificateService) {
	commonName := "my-ca"
	email := "test@mail.com"
	organization := "my-org"

	request := &NewCertificateRequest{
		CommonName:     commonName,
		Email:          email,
		Organization:   organization,
		ExpirationDays: 365,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	subject := cert.Subject
	if subject.CommonName != commonName {
		t.Fatalf("Common name is not correct, expected: [%s] found: [%s] \n", commonName, subject.CommonName)
	}

	if len(subject.Organization) == 0 || subject.Organization[0] != organization {
		t.Fatalf("Organization name is not correct, expected: [%s] found: [%s] \n", organization, subject.Organization)
	}
}

func testNotProvidedCommonName(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		ExpirationDays: 365,
	}
	_, err := (*service).CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: common name")) {
		t.Fatalf("provided not valid email, but error proper error not raised error: [%v]", err)
	}
}

func testUniqueSerialNumber(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 365,
	}
	firstResponse := createCert(t, service, request)
	secondResponse := createCert(t, service, request)

	firstCert := parsePEMToX509Certificate(firstResponse.Certificate)
	secondCert := parsePEMToX509Certificate(secondResponse.Certificate)

	if firstCert.SerialNumber == secondCert.SerialNumber {
		t.Fatalf("Sequentially created two certs' serial numbers are same, should've been different. first: [%s], second: [%s]",
			firstCert.SerialNumber.String(), secondCert.SerialNumber.String())
	}
}

func testExpirationDate(t *testing.T, service *CertificateService) {
	expirationDays := 742
	beforeExpirationDate := time.Now().AddDate(0, 0, expirationDays)

	// ----

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: expirationDays,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	// ---

	certExpireDate := cert.NotAfter
	afterExpirationDate := time.Now().AddDate(0, 0, expirationDays)

	if certExpireDate.After(beforeExpirationDate) && certExpireDate.Before(afterExpirationDate) {
		t.Fatalf("Expiration date is not correct, expected [%d] days later, found: [%s] \n", expirationDays, certExpireDate)
	}
}

func testNotValidExpirationDate(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 0,
	}
	_, err := (*service).CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: expiration days")) {
		t.Fatalf("provided not valid expiration date, must be bigger than 1: [%v]", err)
	}
}

func testNotProvidedExpirationDate(t *testing.T, service *CertificateService) {
	request := &NewCertificateRequest{
		CommonName: "my-ca",
	}

	_, err := (*service).CreateCertificate(request)
	if !(err != nil && strings.Contains(err.Error(), "Validation error: expiration days")) {
		t.Fatalf("provided not valid expiration date, must be bigger than 1: [%v]", err)
	}
}

func testStartDate(t *testing.T, service *CertificateService) {
	beforeCreateCert := time.Now()

	// ----

	request := &NewCertificateRequest{
		CommonName:     "my-ca",
		ExpirationDays: 5,
	}
	response := createCert(t, service, request)
	cert := parsePEMToX509Certificate(response.Certificate)

	// ---

	certExpireDate := cert.NotBefore
	afterCreateCert := time.Now()

	if certExpireDate.After(beforeCreateCert) && cert.NotAfter.Before(afterCreateCert) {
		t.Fatalf("Start date is not correct, should've been now but: [%s] \n", certExpireDate)
	}
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
	if !matched {
		t.Fatalf("expected [%s] but found [%s]", pattern, certificateStr)
	}
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
	if !matched {
		t.Fatalf("expected [%s] but found [%s]", pattern, privateKeyStr)
	}
}

// ---

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

// -----

func createCert(t *testing.T, service *CertificateService, request *NewCertificateRequest) *NewCertificateResponse {
	response, err := (*service).CreateCertificate(request)
	if err != nil {
		t.Fatalf("cert creation failed %v", err)
	}
	return response
}
