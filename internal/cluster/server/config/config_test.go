package config

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/service/factory"
)

func TestParse(t *testing.T) {
	configYaml := `listen-port: 8080
tls-ca-cert: "ca-cert-path"
tls-server-cert: "server-cert-path"
tls-server-cert-key: "server-cert-key-path"
certstore:
  services:
  - name: test-cert-service
    type: Simple`

	config, err := Parse(configYaml)
	assert.NotError(t, err, "parsing failed")

	// ----

	assert.Equal(t, 8080, config.ListenPort)
	assert.Equal(t, "ca-cert-path", config.TlsCACert)
	assert.Equal(t, "server-cert-path", config.TlsServerCert)
	assert.Equal(t, "server-cert-key-path", config.TlsServerCertKey)

	// ----

	certstoreConfig := config.CertStore
	issuerConfigs := certstoreConfig.IssuerConfigs
	assert.Equal(t, 1, len(issuerConfigs))

	assert.Equal(t, "test-cert-service", issuerConfigs[0].Name)
	assert.Equal(t, string(factory.Simple), string(issuerConfigs[0].Type))
}
