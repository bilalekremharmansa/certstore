package config

import (
	certstore_config "bilalekrem.com/certstore/internal/certstore/config"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenPort       int                     `yaml:"listen-port"`
	TlsCACert        string                  `yaml:"tls-ca-cert"`
	TlsServerCert    string                  `yaml:"tls-server-cert"`
	TlsServerCertKey string                  `yaml:"tls-server-cert-key"`
	CertStore        certstore_config.Config `yaml:"certstore"`
}

func Parse(configYaml string) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal([]byte(configYaml), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
