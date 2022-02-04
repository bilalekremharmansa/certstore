package lego

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/x509utils"
	real_lego "github.com/go-acme/lego/v4/lego"
)

func TestNewAdapterWithNewUserRegistration(t *testing.T) {
	email := "certstore@certstore.com"

	dir, err := ioutil.TempDir("/tmp", "test_le_new_adapter_new_user")
	assert.NotError(t, err, "creating temp dir failed")
	defer os.RemoveAll(dir)

	privateKeyPath := dir + "/" + "private-key"
	_, err = NewAdapterWithNewUserRegistration(email, privateKeyPath, nil, real_lego.LEDirectoryStaging)
	assert.NotError(t, err, "creating new lego adapter failed")

	privateKeyContent, err := ioutil.ReadFile(privateKeyPath)
	assert.NotError(t, err, "reading private key failed")

	_, err = x509utils.ParsePemPrivateKey(privateKeyContent)
	assert.NotError(t, err, "decode private key failed")

	// ---

	accountUriPath := privateKeyPath + ".uri"
	accountUriContent, err := ioutil.ReadFile(accountUriPath)
	assert.NotError(t, err, "reading account uri failed")

	accountUri := string(accountUriContent)
	// sample accountUri: https://acme-staging-v02.api.letsencrypt.org/acme/acct/43012568
	assert.True(t, strings.Contains(accountUri, "/acme/acct/"))
}

func TestNewAdapter(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "test_le_new_adapter")
	assert.NotError(t, err, "creating temp dir failed")
	defer os.RemoveAll(dir)

	privateKeyPath := dir + "/" + "private-key"

	email := "certstore@certstore.com"
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NotError(t, err, "generating private key failed")

	encodedPrivateKey := x509utils.EncodePEMPrivateKey(privateKey)
	err = ioutil.WriteFile(privateKeyPath, encodedPrivateKey.Bytes(), 0666)
	assert.NotError(t, err, "writing private key failed")

	accountUriPath := privateKeyPath + ".uri"
	accountUri := real_lego.LEDirectoryStaging + "/acme/acct/11111"
	err = ioutil.WriteFile(accountUriPath, []byte(accountUri), 0666)
	assert.NotError(t, err, "writing account uri key failed")

	// -----

	acmeUser, err := NewAcmeUserWithPrivateKeyFile(email, privateKeyPath)
	assert.NotError(t, err, "creating acme user from with file failed")

	// sample accountUri: https://acme-staging-v02.api.letsencrypt.org/acme/acct/43012568
	assert.True(t, strings.Contains(acmeUser.GetRegistration().URI, "/acme/acct/"))

	_, err = NewAdapter(acmeUser, nil, real_lego.LEDirectoryStaging)
	assert.NotError(t, err, "creating new lego adapter failed")
}
