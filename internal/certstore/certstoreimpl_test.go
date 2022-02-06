package certstore

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	certificate_service "bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/cluster/manager"
	"bilalekrem.com/certstore/internal/testutils"
	"github.com/golang/mock/gomock"
)

func TestCreateCertStoreWithCA(t *testing.T) {
	pemCert := testutils.GetCAPem()
	pemPrivateKey := testutils.GetCAPrivateKey()

	store, err := New([]byte(pemPrivateKey), []byte(pemCert))
	assert.NotError(t, err, "creating cert store failed")

	// --

	assert.NotNil(t, store)
	assert.NotNil(t, store.clusterService)
}

func TestCreateCertStoreWithUnvalidCA(t *testing.T) {
	pemCert := "invalid cert"
	pemPrivateKey := testutils.GetCAPrivateKey()

	_, err := New([]byte(pemPrivateKey), []byte(pemCert))

	assert.Error(t, err, "provided cert is invalid error expected")
}

func TestCreateCertStoreWithUnvalidCAKey(t *testing.T) {
	pemCert := testutils.GetCAPem()
	pemPrivateKey := "invalid key"

	_, err := New([]byte(pemPrivateKey), []byte(pemCert))
	assert.Error(t, err, "provided cert key is invalid error expected")
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

	store := createCertStore(t)
	store.RegisterIssuer("first issuer", firstService)
	store.RegisterIssuer("second issuer", secondService)

	// ----

	store.IssueCertificate("first issuer", firstRequest)
	store.IssueCertificate("second issuer", secondRequest)
}

// -----

func createCertStore(t *testing.T) *certStoreImpl {
	clusterManager, err := manager.NewForInitialization()
	assert.NotError(t, err, "cluster manager could not be created")
	caCert, err := clusterManager.CreateClusterCACertificate("test-cluster")
	assert.NotError(t, err, "cluster ca certificate could not be created")

	// -----

	store, err := New(caCert.PrivateKey, caCert.Certificate)
	assert.NotError(t, err, "cluster certificate service could not be created")

	return store
}
