package certstore

import (
	"testing"

	"bilalekrem.com/certstore/internal/certificate/x509utils"
	"bilalekrem.com/certstore/internal/testutils"
)

func TestCreateCertStore(t *testing.T) {
	store, err := NewWithoutCA()
	if err != nil {
		t.Fatalf("creating cert store failed, %v", err)
	}

	if store == nil {
		t.Fatalf("certstore is nil, should've been created\n")
	}
}

func TestCreateCertStoreWithCA(t *testing.T) {
	pemCert := testutils.GetCAPem()
	pemPrivateKey := testutils.GetCAPrivateKey()

	store, err := New([]byte(pemPrivateKey), []byte(pemCert))
	if err != nil {
		t.Fatalf("creating cert store failed, %v", err)
	}

	// --

	if store == nil {
		t.Fatalf("cerstore is nil, should've been created\n")
	}

	if store.clusterService == nil {
		t.Fatalf("ca key and cert is provided, store.clusterService should not be nil")
	}
}

func TestCreateCertStoreWithUnvalidCA(t *testing.T) {
	pemCert := "invalid cert"
	pemPrivateKey := testutils.GetCAPrivateKey()

	_, err := New([]byte(pemPrivateKey), []byte(pemCert))
	if err == nil {
		t.Fatal("provided cert is invalid, error expected but did not returned")
	}
}

func TestCreateCertStoreWithUnvalidCAKey(t *testing.T) {
	pemCert := testutils.GetCAPem()
	pemPrivateKey := "invalid key"

	_, err := New([]byte(pemPrivateKey), []byte(pemCert))
	if err == nil {
		t.Fatal("provided cert key is invalid, error expected but did not returned")
	}
}

func TestCreateCA(t *testing.T) {
	store, _ := NewWithoutCA()
	clusterName := "my-cluster"
	certificate, err := store.CreateClusterCACertificate(clusterName)
	if err != nil {
		t.Fatalf("cluster ca certificate could not be created, %v", err)
	}

	// -----

	cert, err := x509utils.ParsePemCertificate(certificate.Certificate)
	if err != nil {
		t.Fatalf("generated certificate could not be parsed, %v", err)
	}

	if !cert.IsCA {
		t.Fatal("generated certificate is not CA")
	}

	if cert.Subject.CommonName != cert.Issuer.CommonName {
		t.Fatalf("CA subject and issuer common name are different, expected to be same, subject %s, issuer %s",
			cert.Subject.CommonName, cert.Issuer.CommonName)
	}

	if cert.Subject.CommonName != clusterName {
		t.Fatalf("ca common name is not as expected, expected: [%s], actual: [%s]", clusterName, cert.Subject.CommonName)
	}
}

func TestCreateServerCert(t *testing.T) {
	storeWithoutCA, _ := NewWithoutCA()
	clusterName := "my-cluster"
	caCert, err := storeWithoutCA.CreateClusterCACertificate(clusterName)
	if err != nil {
		t.Fatalf("cluster ca certificate could not be created, %v", err)
	}

	// -----

	store, _ := New(caCert.PrivateKey, caCert.Certificate)
	serverName := "my-server"
	serverCertResponse, err := store.CreateServerCertificate(serverName)
	if err != nil {
		t.Fatalf("server certificate could not be created, %v", err)
	}

	serverCert, err := x509utils.ParsePemCertificate(serverCertResponse.Certificate)

	if serverCert.Subject.CommonName != serverName {
		t.Fatalf("server cert common name is not as expected, expected: [%s], actual: [%s]", serverName, serverCert.Subject.CommonName)
	}

	if serverCert.IsCA {
		t.Fatal("generated server certificate is CA, should've been server")
	}
}

func TestCreateWorkerCert(t *testing.T) {
	storeWithoutCA, _ := NewWithoutCA()
	clusterName := "my-cluster"
	caCert, err := storeWithoutCA.CreateClusterCACertificate(clusterName)
	if err != nil {
		t.Fatalf("cluster ca certificate could not be created, %v", err)
	}

	// -----

	store, _ := New(caCert.PrivateKey, caCert.Certificate)
	workerName := "my-worker"
	workerCertResponse, err := store.CreateWorkerCertificate(workerName)
	if err != nil {
		t.Fatalf("worker certificate could not be created, %v", err)
	}

	workerCert, err := x509utils.ParsePemCertificate(workerCertResponse.Certificate)

	if workerCert.Subject.CommonName != workerName {
		t.Fatalf("worker cert common name is not as expected, expected: [%s], actual: [%s]", workerName, workerCert.Subject.CommonName)
	}

	if workerCert.IsCA {
		t.Fatal("generated worker certificate is CA, should've been server")
	}
}
