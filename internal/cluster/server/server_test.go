package server

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/cluster/server/config"
)

func TestValidateConfig(t *testing.T) {
	conf := getConfig()
	err := validateConfig(conf)
	assert.NotError(t, err, "validation failed")
}

func TestValidateConfigMissingListenPort(t *testing.T) {
	conf := getConfig()
	conf.ListenPort = 0
	err := validateConfig(conf)
	assert.Error(t, err, "validation failed: missing listen port")
}

func TestValidateConfigMissingTlsCACert(t *testing.T) {
	conf := getConfig()
	conf.TlsCACert = ""
	err := validateConfig(conf)
	assert.Error(t, err, "validation failed: missing ca cert")
}

func TestValidateConfigMissingTlsServerCert(t *testing.T) {
	conf := getConfig()
	conf.TlsServerCert = ""
	err := validateConfig(conf)
	assert.Error(t, err, "validation failed: missing server cert")
}

func TestValidateConfigMissingTlsServerCertKey(t *testing.T) {
	conf := getConfig()
	conf.TlsServerCertKey = ""
	err := validateConfig(conf)
	assert.Error(t, err, "validation failed: missing server cert key")
}

func getConfig() *config.Config {
	conf := &config.Config{}
	conf.ListenPort = 10000
	conf.TlsCACert = "tls-ca-cert"
	conf.TlsServerCert = "tls-server-cert"
	conf.TlsServerCertKey = "tls"
	return conf
}