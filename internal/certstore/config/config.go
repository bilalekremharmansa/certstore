package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	service_factory "bilalekrem.com/certstore/internal/certificate/service/factory"
	"bilalekrem.com/certstore/internal/logging"
)

type Config struct {
	IssuerConfigs []CertificateServiceConfig      `yaml:"services"`
}

type CertificateServiceConfig struct {
	Name string                      `yaml:"name"`
	Type service_factory.ServiceType `yaml:"type"`
	Args map[string]string           `yaml:"args"`
}

// ------

func ParseFile(path string) (*Config, error) {
	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parse(yamlBytes)
}

func ParseYaml(yaml string) (*Config, error) {
	yamlBytes := []byte(yaml)
	return parse(yamlBytes)
}

func parse(yamlBytes []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(yamlBytes, config)
	if err != nil {
		logging.GetLogger().Errorf("Parsing config yaml failed %s, %v", string(yamlBytes), err)
		return nil, err
	}

	// ----

	err = validate(config)
	if err != nil {
		logging.GetLogger().Errorf("Validating config failed %v\n", err)
		return nil, err
	}

	return config, nil
}

func validate(config *Config) error {
	for _, issuerConfig := range config.IssuerConfigs {
		if issuerConfig.Name == "" {
			return errors.New("issuer config name is empty, 'name' is required")
		}

		if issuerConfig.Type != service_factory.Simple &&
			issuerConfig.Type != service_factory.CertificateAuthority &&
			issuerConfig.Type != service_factory.LetsEncrypt {
			return errors.New(fmt.Sprintf("issuer config service type is unknown, 'ServiceType' is required, %s",
				string(issuerConfig.Type)))
		}
	}
	return nil
}
