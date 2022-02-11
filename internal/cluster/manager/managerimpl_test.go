package manager

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/x509utils"
)

func TestCreateCA(t *testing.T) {
	clusterManager, err := NewForInitialization()
	assert.NotError(t, err, "cluster manager could not be created")
	clusterName := "my-cluster"
	certificate, err := clusterManager.CreateClusterCACertificate(clusterName)
	assert.NotError(t, err, "cluster ca certificate could not be created")

	// -----

	cert, err := x509utils.ParsePemCertificate(certificate.Certificate)
	assert.NotError(t, err, "certificate could not be parsed")

	assert.True(t, cert.IsCA)
	assert.EqualM(t, cert.Subject.CommonName, cert.Issuer.CommonName, "CA subject and issuer common name are different")
	assert.Equal(t, clusterName, cert.Subject.CommonName)
}

func TestCreateServerCert(t *testing.T) {
	clusterManager := createClusterManagerWithCA(t)
	serverName := "my-server"
	serverCertResponse, err := clusterManager.CreateServerCertificate(serverName)
	assert.NotError(t, err, "server certificate could not be created")

	serverCert, err := x509utils.ParsePemCertificate(serverCertResponse.Certificate)

	assert.Equal(t, serverName, serverCert.Subject.CommonName)
	assert.False(t, serverCert.IsCA)
}

func TestCreateAgentCert(t *testing.T) {
	clusterManager := createClusterManagerWithCA(t)
	agentName := "my-agent"
	agentCertResponse, err := clusterManager.CreateAgentCertificate(agentName)
	assert.NotError(t, err, "agent certificate could not be created")

	agentCert, err := x509utils.ParsePemCertificate(agentCertResponse.Certificate)

	assert.Equal(t, agentName, agentCert.Subject.CommonName)
	assert.False(t, agentCert.IsCA)
}

func createClusterManagerWithCA(t *testing.T) ClusterManager {
	initialClusterManager, err := NewForInitialization()
	assert.NotError(t, err, "cluster manager could not be created")

	response, err := initialClusterManager.CreateClusterCACertificate("test-cluster")
	assert.NotError(t, err, "ca cert could not be created")

	clusterManager, err := NewFromCA(response.Certificate, response.PrivateKey)
	assert.NotError(t, err, "ca cert could not be created")
	return clusterManager
}
