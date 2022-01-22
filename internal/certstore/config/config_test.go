package config

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	service_factory "bilalekrem.com/certstore/internal/certificate/service/factory"
)

func TestParseConfig(t *testing.T) {
	config, err := ParseYaml(`cluster:
  private-key: cluster-private-key-file-path
  certificate: cluster-certificate-file-path
services:
  - name: test-cert-service
    type: Simple
    args:
      private-key: simple-private-key-file-path
      certificate: simple-certificate-file-path`)

	assert.NotError(t, err, "error occurred while parsing yaml")

	clusterConfig := config.ClusterConfig
	assert.Equal(t, "cluster-certificate-file-path", clusterConfig.CertificatePath)
	assert.Equal(t, "cluster-private-key-file-path", clusterConfig.PrivateKeyPath)

	// -----

	issuerConfigs := config.IssuerConfigs
	assert.Equal(t, 1, len(issuerConfigs))

	issuerConfig := issuerConfigs[0]
	assert.Equal(t, "test-cert-service", issuerConfig.Name)

	assert.Equal(t, service_factory.Simple, issuerConfig.Type)

	issuerConfigArgs := issuerConfig.Args

	assert.Equal(t, "simple-certificate-file-path", issuerConfigArgs["certificate"])
	assert.Equal(t, "simple-private-key-file-path", issuerConfigArgs["private-key"])
}

func TestIssuerServiceNameEmpty(t *testing.T) {
	_, err := ParseYaml(`services:
  - type: Simple
    args:
      private-key: simple-private-key-file-path
      certificate: simple-certificate-file-path`)

	assert.Error(t, err, "issuer config name is empty")
}

func TestIssuerServiceTypeEmpty(t *testing.T) {
	_, err := ParseYaml(`services:
  - name: test-cert-service`)

	assert.Error(t, err, "issuer config type is empty")
}

func TestIssuerServiceTypeSimple(t *testing.T) {
	config, err := ParseYaml(`services:
  - name: test-cert-service
    type: Simple`)

	assert.NotError(t, err, "parsing yaml failed")
	assert.Equal(t, service_factory.Simple, config.IssuerConfigs[0].Type)
}

func TestIssuerServiceTypeCertificateAuthority(t *testing.T) {
	config, err := ParseYaml(`services:
  - name: test-cert-service
    type: CertificateAuthority`)

	assert.NotError(t, err, "parsing yaml failed")
	assert.DeepEqual(t, service_factory.CertificateAuthority, string(config.IssuerConfigs[0].Type))
}
