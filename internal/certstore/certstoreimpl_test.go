package certstore

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	certificate_service "bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certstore/config"
	"github.com/golang/mock/gomock"
)

func TestCreateCertStoreWithConfig(t *testing.T) {
	store := createWithConfig(t)
	assert.NotNil(t, store)

	assert.Equal(t, 1, len(store.certIssuers))
}

func TestIssueCertificate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	firstRequest := &certificate_service.NewCertificateRequest{CommonName: "first cert"}
	firstService := certificate_service.NewMockCertificateService(ctrl)
	firstService.
		EXPECT().
		CreateCertificate(gomock.Eq(firstRequest)).
		MinTimes(1)

	secondRequest := &certificate_service.NewCertificateRequest{CommonName: "second cert"}
	secondService := certificate_service.NewMockCertificateService(ctrl)
	secondService.
		EXPECT().
		CreateCertificate(gomock.Eq(secondRequest)).
		MinTimes(1)

	// ----

	store := createWithConfig(t)
	store.RegisterIssuer("first issuer", firstService)
	store.RegisterIssuer("second issuer", secondService)

	// ----

	store.IssueCertificate("first issuer", firstRequest)
	store.IssueCertificate("second issuer", secondRequest)
}

// -----

func createWithConfig(t *testing.T) *certStoreImpl {
	configYaml := `services:
  - name: test-cert-service
    type: Simple
    args:
      private-key: simple-private-key-file-path
      certificate: simple-certificate-file-path`
	conf, err := config.ParseYaml(configYaml)
	assert.NotError(t, err, "parsing certstore config failed")

	store, err := NewFromConfig(conf)
	assert.NotError(t, err, "parsing certstore config failed")

	return store
}
