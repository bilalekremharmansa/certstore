package config

import (
	"bilalekrem.com/certstore/internal/pipeline"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddr       string `yaml:"server-address"`
	TlsCACert        string `yaml:"tls-ca-cert"`
	TlsWorkerCert    string `yaml:"tls-worker-cert"`
	TlsWorkerCertKey string `yaml:"tls-worker-cert-key"`

	Pipelines []pipeline.PipelineConfig `yaml:"pipelines"`
	Jobs      []JobConfig               `yaml:"jobs"`
}

type JobConfig struct {
	Name     string `yaml:"name"`
	Pipeline string `yaml:"pipeline"`
}

func Parse(configYaml string) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal([]byte(configYaml), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
