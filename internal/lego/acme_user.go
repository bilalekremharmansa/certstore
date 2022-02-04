package lego

import (
	"crypto"
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/go-acme/lego/v4/registration"
)

type AcmeUser struct {
	email        string
	registration *registration.Resource
	key          crypto.PrivateKey
}

func NewAcmeUser(email string, key crypto.PrivateKey) (*AcmeUser, error) {
	return &AcmeUser{email: email, key: key}, nil
}

func NewAcmeUserWithPrivateKeyFile(email string, privateKeyPath string) (*AcmeUser, error) {
	privateKeyContent, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		logging.GetLogger().Errorf("reading private key failed %v", err)
		return nil, err
	}

	privateKey, err := x509utils.ParsePemPrivateKey(privateKeyContent)
	if err != nil {
		logging.GetLogger().Errorf("decoding pem failed for acme user %v", err)
		return nil, err
	}

	// ---

	accountUriPath := privateKeyPath + ".uri"
	accountUriContent, err := ioutil.ReadFile(accountUriPath)
	if err != nil {
		logging.GetLogger().Errorf("reading account uri failed %v", err)
		return nil, err
	}

	reg := &registration.Resource{URI: string(accountUriContent)}
	return &AcmeUser{email: email, key: privateKey, registration: reg}, nil
}

func (u *AcmeUser) GetEmail() string {
	return u.email
}
func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
