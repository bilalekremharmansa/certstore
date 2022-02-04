package lego

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/certificate"
	"github.com/go-acme/lego/challenge"
	real_lego "github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
)

type legoAdapterImpl struct {
	legoClient *real_lego.Client
}

func NewAdapter(user *AcmeUser, provider challenge.Provider, caDirUrl string) (*legoAdapterImpl, error) {
	config := real_lego.NewConfig(user)
	config.CADirURL = caDirUrl
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := real_lego.NewClient(config)
	if err != nil {
		logging.GetLogger().Errorf("creating lego client failed %v", err)
		return nil, err
	}

	// -----

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		logging.GetLogger().Errorf("setting new dns 01 provider failed %v", err)
		return nil, err
	}

	// -----

	return &legoAdapterImpl{legoClient: client}, nil
}

// this function will generate a new lets encrypt user and will save private key to 'userPrivateKeyPath'
func NewAdapterWithNewUserRegistration(userEmail string, userPrivateKeyPath string, provider challenge.Provider, caDirUrl string) (*legoAdapterImpl, error) {
	user, err := createAndRegisterNewUser(userEmail, userPrivateKeyPath, caDirUrl)
	if err != nil {
		return nil, err
	}

	return NewAdapter(user, provider, caDirUrl)
}

func (c *legoAdapterImpl) Obtain(req certificate.ObtainRequest) (*certificate.Resource, error) {
	certificates, err := c.legoClient.Certificate.Obtain(req)
	if err != nil {
		return nil, err
	}

	return certificates, nil
}

// --------

func createAndRegisterNewUser(email string, userPrivateKeyPath string, caDirUrl string) (*AcmeUser, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logging.GetLogger().Debugf("generating ca private key failed: [%v]", err)
		return nil, err
	}

	user, err := NewAcmeUser(email, privateKey)
	if err != nil {
		logging.GetLogger().Errorf("generating user key failed %v", err)
		return nil, err
	}

	// -------

	config := real_lego.NewConfig(user)
	config.CADirURL = caDirUrl
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := real_lego.NewClient(config)
	if err != nil {
		logging.GetLogger().Errorf("creating lego client failed %v", err)
		return nil, err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.registration = reg

	// -----

	encodedPrivateKey := x509utils.EncodePEMPrivateKey(privateKey)
	err = ioutil.WriteFile(userPrivateKeyPath, encodedPrivateKey.Bytes(), 0666)
	if err != nil {
		logging.GetLogger().Errorf("write generate user private key to file failed %v", err)
		return nil, err
	}

	return user, nil
}
