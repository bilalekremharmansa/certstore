package config

import (
	"bilalekrem.com/certstore/internal/pipeline"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cluster   ClusterConfig             `yaml:"cluster"`
	Pipelines []pipeline.PipelineConfig `yaml:"pipelines"`
}

type ClusterConfig struct {
	ServerAddr       string `yaml:"server-address"`
	TlsCACert        string `yaml:"tls-ca-cert"`
	TlsWorkerCert    string `yaml:"tls-worker-cert"`
	TlsWorkerCertKey string `yaml:"tls-worker-cert-key"`
}

func Parse(configYaml string) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal([]byte(configYaml), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
