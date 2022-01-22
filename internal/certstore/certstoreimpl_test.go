package certstore

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	certificate_service "bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/testutils"
	"github.com/golang/mock/gomock"
)

func TestCreateCertStore(t *testing.T) {
	store, err := NewWithoutCA()
	assert.NotError(t, err, "creating cert store failed")

	assert.NotNil(t, store)
}

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

func TestCreateCA(t *testing.T) {
	store, _ := NewWithoutCA()
	clusterName := "my-cluster"
	certificate, err := store.CreateClusterCACertificate(clusterName)
	assert.NotError(t, err, "cluster ca certificate could not be created")

	// -----

	cert, err := x509utils.ParsePemCertificate(certificate.Certificate)
	assert.NotError(t, err, "certificate could not be parsed")

	assert.True(t, cert.IsCA)
	assert.EqualM(t, cert.Subject.CommonName, cert.Issuer.CommonName, "CA subject and issuer common name are different")
	assert.Equal(t, clusterName, cert.Subject.CommonName)
}

func TestCreateServerCert(t *testing.T) {
	store := createCertStore(t)
	serverName := "my-server"
	serverCertResponse, err := store.CreateServerCertificate(serverName)
	assert.NotError(t, err, "server certificate could not be created")

	serverCert, err := x509utils.ParsePemCertificate(serverCertResponse.Certificate)

	assert.Equal(t, serverName, serverCert.Subject.CommonName)
	assert.False(t, serverCert.IsCA)
}

func TestCreateWorkerCert(t *testing.T) {
	store := createCertStore(t)
	workerName := "my-worker"
	workerCertResponse, err := store.CreateWorkerCertificate(workerName)
	assert.NotError(t, err, "worker certificate could not be created")

	workerCert, err := x509utils.ParsePemCertificate(workerCertResponse.Certificate)

	assert.Equal(t, workerName, workerCert.Subject.CommonName)
	assert.False(t, workerCert.IsCA)
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

	store, _ := NewWithoutCA()
	store.RegisterIssuer("first issuer", firstService)
	store.RegisterIssuer("second issuer", secondService)

	// ----

	store.IssueCertificate("first issuer", firstRequest)
	store.IssueCertificate("second issuer", secondRequest)
}

// -----

func createCertStore(t *testing.T) *certStoreImpl {
	storeWithoutCA, _ := NewWithoutCA()
	clusterName := "my-cluster"
	caCert, err := storeWithoutCA.CreateClusterCACertificate(clusterName)
	assert.NotError(t, err, "cluster ca certificate could not be created")

	// -----

	store, err := New(caCert.PrivateKey, caCert.Certificate)
	assert.NotError(t, err, "cluster certificate service could not be created")

	return store
}
