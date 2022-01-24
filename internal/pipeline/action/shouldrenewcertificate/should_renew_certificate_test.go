package shouldrenewcertificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/x509utils"
)

func TestShouldNotRenewCertificate(t *testing.T) {
	daysToExpire := 100
	cert := &x509.Certificate{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, daysToExpire),
	}

	shouldRenew := shouldRenewCertificate(cert)
	assert.False(t, shouldRenew)
}

func TestShouldRenewCertificateAlreadyExpired(t *testing.T) {
	cert := &x509.Certificate{
		NotBefore: time.Now().AddDate(0, 0, -10),
		NotAfter:  time.Now().AddDate(0, 0, -8),
	}

	shouldRenew := shouldRenewCertificate(cert)
	assert.True(t, shouldRenew)
}

func TestShouldRenewCertificate(t *testing.T) {
	daysToExpire := 3
	cert := &x509.Certificate{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 0, daysToExpire),
	}

	shouldRenew := shouldRenewCertificate(cert)
	assert.True(t, shouldRenew)
}

func TestRequiredArgumentCertificatePath(t *testing.T) {
	args := make(map[string]string)

	err := NewShouldRenewCertificateAction().Run(nil, args)
	assert.ErrorContains(t, err, "required argument")
}

func TestRun(t *testing.T) {
	testRun(t, 100, false)
	testRun(t, 50, false)
	testRun(t, 26, false)

	testRun(t, 20, true)
	testRun(t, 25, true)
	testRun(t, -5, true)
}

func TestCertificateFileNotFound(t *testing.T) {
	args := make(map[string]string)
	args[ARGS_CERTIFICATE_PATH] = "./not-exist-cert-file-path"

	err := NewShouldRenewCertificateAction().Run(nil, args)
	assert.NotError(t, err, "file not found, that means should be renewed")
}

func TestCertificateFoundButNotInCorrectFormat(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_should_renew_certificate_action_wrong_file_format")
	assert.NotError(t, err, "creating temp dir")
	defer os.RemoveAll(dir)

	// ----

	certificatePath := fmt.Sprintf("%s/test.crt", dir)
	err = ioutil.WriteFile(certificatePath, []byte("test certificate content - but not in PEM format"), 0666)

	args := make(map[string]string)
	args[ARGS_CERTIFICATE_PATH] = certificatePath

	err = NewShouldRenewCertificateAction().Run(nil, args)
	assert.NotError(t, err, "file content is not in PEM format, that means should be renewed")
}

func testRun(t *testing.T, expirationDate int, shouldRenew bool) {
	dir, err := ioutil.TempDir("/tmp", "test_should_renew_certificate_action")
	assert.NotError(t, err, "creating temp dir")
	defer os.RemoveAll(dir)

	// ----

	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, expirationDate),
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NotError(t, err, "generating ca private key failed")

	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &privateKey.PublicKey, privateKey)
	assert.NotError(t, err, "creating certificate failed")

	// ----

	certificatePath := fmt.Sprintf("%s/test.crt", dir)
	certificatePem := x509utils.EncodePEMCert(certBytes)

	err = ioutil.WriteFile(certificatePath, certificatePem.Bytes(), 0666)
	assert.NotError(t, err, "writing certificate failed")

	// -----

	args := make(map[string]string)
	args[ARGS_CERTIFICATE_PATH] = fmt.Sprintf("%s/test.crt", dir)

	err = NewShouldRenewCertificateAction().Run(nil, args)
	if shouldRenew {
		assert.NotError(t, err, "certificate should bew renewed, but decided not to")
	} else {
		assert.Error(t, err, "no need to renew is expected, but decided to renew")
	}

}
