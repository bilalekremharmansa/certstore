package config

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
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
      - name: second-action
jobs:
  - name: "first-pipeline job"
    pipeline: "first-pipeline"`

	config, err := Parse(configYaml)
	assert.NotError(t, err, "parsing failed")

	// ----

	clusterConfig := config.Cluster
	assert.Equal(t, "addr:port", clusterConfig.ServerAddr)
	assert.Equal(t, "ca-cert-path", clusterConfig.TlsCACert)
	assert.Equal(t, "worker-cert-path", clusterConfig.TlsWorkerCert)
	assert.Equal(t, "worker-cert-key-path", clusterConfig.TlsWorkerCertKey)

	// -----

	pipelines := config.Pipelines
	assert.Equal(t, 2, len(pipelines))
	assert.Equal(t, "first-pipeline", pipelines[0].Name)
	assert.Equal(t, 2, len(pipelines[0].Actions))

	// -----

	jobs := config.Jobs
	assert.Equal(t, 1, len(jobs))
	assert.Equal(t, "first-pipeline job", jobs[0].Name)
	assert.Equal(t, "first-pipeline", jobs[0].Pipeline)
}
