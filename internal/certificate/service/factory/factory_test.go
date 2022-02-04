package factory

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/testutils"
)

func TestNewSimpleCertificateService(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_new_cert_service")
	assert.NotError(t, err, "creating temp dir failed")
	defer os.RemoveAll(dir)

	// ------

	privateKeyPath := fmt.Sprintf("%s/ca.key", dir)
	privateKey := testutils.GetCAPrivateKey()
	ioutil.WriteFile(privateKeyPath, []byte(privateKey), 0666)

	certPath := fmt.Sprintf("%s/ca.crt", dir)
	certPem := testutils.GetCAPem()
	ioutil.WriteFile(certPath, []byte(certPem), 0666)

	// -----

	args := make(map[string]string)
	args["private-key"] = privateKeyPath
	args["certificate"] = certPath

	service := NewService(Simple, args)
	assert.NotNil(t, service)
}

func TestCACertificateService(t *testing.T) {
	service := NewService(CertificateAuthority, nil)
	assert.NotNil(t, service)
}

func TestLetsEncrypt(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_new_cert_service")
	assert.NotError(t, err, "creating temp dir failed")
	defer os.RemoveAll(dir)

	// -----

	privateKeyPath := fmt.Sprintf("%s/user.key", dir)

	args := make(map[string]string)
	args["private-key"] = privateKeyPath
	args["email"] = "test@certstore.com"
	args["provider"] = "mock"

	service := NewService(LetsEncrypt, args)
	assert.NotNil(t, service)
}

func TestLetsEncryptMissingFields(t *testing.T) {
	// email is missing
	args := make(map[string]string)
	args["private-key"] = "private key path"
	args["provider"] = "mock"

	service := NewService(LetsEncrypt, nil)
	assert.Nil(t, service)

	// -----

	// private key is missing
	args = make(map[string]string)
	args["email"] = "test@certstore.com"
	args["provider"] = "mock"

	service = NewService(LetsEncrypt, nil)
	assert.Nil(t, service)

	// -----

	// provider is missing
	args = make(map[string]string)
	args["private-key"] = "private key path"
	args["email"] = "test@certstore.com"

	service = NewService(LetsEncrypt, nil)
	assert.Nil(t, service)
}

func TestUnknownServiceShouldBeNil(t *testing.T) {
	service := NewService(Unknown, nil)
	assert.Nil(t, service)
}

func TestUnrelatedServiceShouldBeNil(t *testing.T) {
	service := NewService("test", nil)
	assert.Nil(t, service)
}
