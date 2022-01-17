package config

import (
	"testing"

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

	if err != nil {
		t.Fatalf("error occurred while parsing yaml, %v", err)
	}

	clusterConfig := config.ClusterConfig
	if clusterConfig.PrivateKeyPath != "cluster-private-key-file-path" {
		t.Fatalf("cluster config private key path is not matching")
	}
	if clusterConfig.CertificatePath != "cluster-certificate-file-path" {
		t.Fatalf("cluster config certificate path is not matching")
	}

	// -----

	issuerConfigs := config.IssuerConfigs
	if len(issuerConfigs) != 1 {
		t.Fatalf("Unexpected issuer config")
	}

	issuerConfig := issuerConfigs[0]
	if issuerConfig.Name != "test-cert-service" {
		t.Fatalf("issuer config service name is not matching")
	}
	if issuerConfig.Type != service_factory.Simple {
		t.Fatalf("issuer config service type is not matching")
	}

	issuerConfigArgs := issuerConfig.Args
	if issuerConfigArgs["private-key"] != "simple-private-key-file-path" {
		t.Fatalf("issuer config arg privateKey is not matching")
	}
	if issuerConfigArgs["certificate"] != "simple-certificate-file-path" {
		t.Fatalf("issuer config arg certificate is not matching")
	}
}

func TestIssuerServiceNameEmpty(t *testing.T) {
	_, err := ParseYaml(`services:
  - type: Simple
    args:
      private-key: simple-private-key-file-path
      certificate: simple-certificate-file-path`)

	if err == nil {
		t.Fatal("error is expected since issuer config name is empty, but not found")
	}
}

func TestIssuerServiceTypeEmpty(t *testing.T) {
	_, err := ParseYaml(`services:
  - name: test-cert-service`)

	if err == nil {
		t.Fatal("error is expected since issuer config type is empty, but not found")
	}
}

func TestIssuerServiceTypeSimple(t *testing.T) {
	config, err := ParseYaml(`services:
  - name: test-cert-service
    type: Simple`)

	if err != nil {
		t.Fatalf("error occurred while parsing yaml, %v", err)
	}

	if config.IssuerConfigs[0].Type != service_factory.Simple {
		t.Fatalf("issuer config service type must be simple, found %s", config.IssuerConfigs[0].Type)
	}
}

func TestIssuerServiceTypeCertificateAuthority(t *testing.T) {
	config, err := ParseYaml(`services:
  - name: test-cert-service
    type: CertificateAuthority`)

	if err != nil {
		t.Fatalf("error occurred while parsing yaml, %v", err)
	}

	if config.IssuerConfigs[0].Type != service_factory.CertificateAuthority {
		t.Fatalf("issuer config service type must be CertificateAuthority, found %s", config.IssuerConfigs[0].Type)
	}
}
