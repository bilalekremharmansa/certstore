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

func TestUnknownServiceShouldBeNil(t *testing.T) {
	service := NewService(Unknown, nil)
	assert.Nil(t, service)
}

func TestUnrelatedServiceShouldBeNil(t *testing.T) {
	service := NewService("test", nil)
	assert.Nil(t, service)
}
