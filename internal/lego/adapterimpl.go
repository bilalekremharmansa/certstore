package lego

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"

	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/logging"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	real_lego "github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
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

	err = client.Challenge.SetDNS01Provider(provider, dns01.DisableCompletePropagationRequirement())
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
		logging.GetLogger().Errorf("Obtaining certificate failed request:%v, %v", req, err)
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

	// ----

	err = saveNewUserAssets(user, userPrivateKeyPath)
	if err != nil {
		logging.GetLogger().Errorf("saving acme user assets %v", err)
		return nil, err
	}

	return user, nil
}

func saveNewUserAssets(user *AcmeUser, userPrivateKeyPath string) error {
	privateKey := user.GetPrivateKey().(*rsa.PrivateKey)
	encodedPrivateKey := x509utils.EncodePEMPrivateKey(privateKey)
	err := ioutil.WriteFile(userPrivateKeyPath, encodedPrivateKey.Bytes(), 0666)
	if err != nil {
		logging.GetLogger().Errorf("write generate user private key to file failed %v", err)
		return err
	}

	accountUriPath := userPrivateKeyPath + ".uri"
	reg := user.GetRegistration()
	err = ioutil.WriteFile(accountUriPath, []byte(reg.URI), 0666)
	if err != nil {
		logging.GetLogger().Errorf("write account uri to file failed %v", err)
		return err
	}

	return nil
}
