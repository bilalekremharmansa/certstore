package factory

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"bilalekrem.com/certstore/internal/testutils"
)

func TestNewSimpleCertificateService(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_new_cert_service")
	if err != nil {
		t.Fatalf("error occurred while creating temp dir, %v", err)
	}
	defer os.RemoveAll(dir)

	// ------

	privateKeyPath := fmt.Sprintf("%s/ca.key", dir)
	privateKey := testutils.GetCAPrivateKey()
	ioutil.WriteFile(privateKeyPath, []byte(privateKey), 0666)

	certPath := fmt.Sprintf("%s/ca.crt", dir)
	certPem := testutils.GetCAPem()
	ioutil.WriteFile(certPath, []byte(certPem), 0666)

	// -----

	args := make(map[string]interface{})
	args["private-key"] = privateKeyPath
	args["certificate"] = certPath

	service := NewService(Simple, args)
	if service == nil {
		t.Fatal("error occurred while creating new simple service")
	}
}

func TestCACertificateService(t *testing.T) {
	service := NewService(CertificateAuthority, nil)
	if service == nil {
		t.Fatal("error occurred while creating new ca certificate service")
	}
}
