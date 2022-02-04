package letsencrypt

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/lego"
	"github.com/go-acme/lego/certificate"
	"github.com/golang/mock/gomock"
)

func TestProvider(t *testing.T) {
	_, err := getProvider("mock")
	assert.NotError(t, err, "mock dns provider not found")
}

func TestCreateCertificate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := lego.NewMockLegoAdapter(ctrl)
	leService := &letsEncryptCertificateService{lego: adapter}

	commonName := "certstore.com"
	responseCert := []byte("test certificate content")
	responsePrivateKey := []byte("test private key content")

	adapter.
		EXPECT().
		Obtain(gomock.Any()).
		DoAndReturn(func(req certificate.ObtainRequest) (*certificate.Resource, error) {
			assert.Equal(t, 1, len(req.Domains))
			assert.Equal(t, req.Domains[0], commonName)

			return &certificate.Resource{
				Certificate: responseCert,
				PrivateKey:  responsePrivateKey,
			}, nil
		})

	request := &service.NewCertificateRequest{
		CommonName: commonName,
	}
	response, err := leService.CreateCertificate(request)
	assert.NotNil(t, response)
	assert.NotError(t, err, "creating lets encrypt cert failed")

	assert.DeepEqual(t, responseCert, response.Certificate)
	assert.DeepEqual(t, responsePrivateKey, response.PrivateKey)
}

func TestCreateCertificateWithSans(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := lego.NewMockLegoAdapter(ctrl)
	leService := &letsEncryptCertificateService{lego: adapter}

	commonName := "certstore.com"
	sans := []string{"test.certstore.com", "live.certstore.com"}
	responseCert := []byte("test certificate content")
	responsePrivateKey := []byte("test private key content")

	adapter.
		EXPECT().
		Obtain(gomock.Any()).
		DoAndReturn(func(req certificate.ObtainRequest) (*certificate.Resource, error) {
			assert.Equal(t, 3, len(req.Domains))
			assert.Equal(t, req.Domains[0], commonName)
			assert.Equal(t, req.Domains[1], sans[0])
			assert.Equal(t, req.Domains[2], sans[1])

			return &certificate.Resource{
				Certificate: responseCert,
				PrivateKey:  responsePrivateKey,
			}, nil
		})

	request := &service.NewCertificateRequest{
		CommonName:              commonName,
		SubjectAlternativeNames: sans,
	}
	response, err := leService.CreateCertificate(request)
	assert.NotNil(t, response)
	assert.NotError(t, err, "creating lets encrypt cert failed")

	assert.DeepEqual(t, responseCert, response.Certificate)
	assert.DeepEqual(t, responsePrivateKey, response.PrivateKey)
}

func TestCreateCertificateMissingCommonName(t *testing.T) {
	leService := &letsEncryptCertificateService{lego: nil}

	request := &service.NewCertificateRequest{
		SubjectAlternativeNames: []string{},
	}
	_, err := leService.CreateCertificate(request)
	assert.ErrorContains(t, err, "Validation error: common name")
}
