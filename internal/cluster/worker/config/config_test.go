package config

import (
	"testing"
)

func TestParse(t *testing.T) {
	configYaml := `cluster:
  server-address: "addr:port"
  tls-ca-cert: "ca-cert-path"
  tls-worker-cert: "worker-cert-path"
  tls-worker-cert-key: "worker-cert-key-path"
pipelines:
  - name: first-pipeline
    actions:
      - name: first-action
        args:
          command: first-arg
      - name: first-action
        args:
          my-arg: my-val
  - name: second-pipeline
    actions:
      - name: second-action`

	config, err := Parse(configYaml)
	if err != nil {
		t.Fatalf("parsing failed: %v\n", err)
	}

	// ----

	clusterConfig := config.Cluster
	if clusterConfig.ServerAddr != "addr:port" {
		t.Fatalf("server address is not correct, found: %s", clusterConfig.ServerAddr)
	}

	if clusterConfig.TlsCACert != "ca-cert-path" {
		t.Fatalf("ca cert is not correct, found: %s", clusterConfig.TlsCACert)
	}

	if clusterConfig.TlsWorkerCert != "worker-cert-path" {
		t.Fatalf("worker cert is not correct, found: %s", clusterConfig.TlsWorkerCert)
	}
	if clusterConfig.TlsWorkerCertKey != "worker-cert-key-path" {
		t.Fatalf("worker cert key is not correct, found: %s", clusterConfig.TlsWorkerCertKey)
	}

	// -----

	pipelines := config.Pipelines
	if len(pipelines) != 2 {
		t.Fatalf("size of pipeline is not correct, found: %d", len(pipelines))
	}

	if pipelines[0].Name != "first-pipeline" {
		t.Fatalf("Pipeline name is not correct, found: %s", pipelines[0].Name)
	}

	if len(pipelines[0].Actions) != 2 {
		t.Fatalf("size of first pipeline action is not correct, found: %d", len(pipelines[0].Actions))
	}

}
