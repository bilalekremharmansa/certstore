package worker

import (
	"io/ioutil"
	"os"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/cluster/manager"
	"bilalekrem.com/certstore/internal/cluster/worker/config"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/store"
)

func TestNewFromFile(t *testing.T) {
	testNewFromFile(t, false)
}

func TestNewFromFileSkipJobInitialization(t *testing.T) {
	testNewFromFile(t, true)
}

func TestInitPipelineConfigSuccess(t *testing.T) {
	worker := &Worker{
		pipelineStore: store.New(),
	}

	actionStore := action.NewActionStore()
	actionStore.Put("action-one", &action.MockAction{})
	actionStore.Put("action-two", &action.MockAction{})

	pipelineConfigs := []pipeline.PipelineConfig{
		{Name: "test-pipeline",
			Actions: []pipeline.PipelineActionConfig{
				{Name: "action-one"},
				{Name: "action-two"},
			}},
	}
	err := worker.initPipelines(pipelineConfigs, actionStore)
	assert.NotError(t, err, "pipeline initialization failed")
}

func TestInitPipelineConfigFail(t *testing.T) {
	worker := &Worker{
		pipelineStore: store.New(),
	}

	actionStore := action.NewActionStore()
	actionStore.Put("action-one", &action.MockAction{})

	pipelineConfigs := []pipeline.PipelineConfig{
		{Name: "test-pipeline",
			Actions: []pipeline.PipelineActionConfig{
				{Name: "action-one"},
				{Name: "action-two"},
			}},
	}
	err := worker.initPipelines(pipelineConfigs, actionStore)
	assert.Error(t, err, "should've been failed, action-two is missing in store")
}

func TestInitJobConfigSuccess(t *testing.T) {
	worker := &Worker{
		pipelineStore: store.New(),
	}

	pipelineConfigs := []pipeline.PipelineConfig{
		{Name: "test-pipeline",
			Actions: []pipeline.PipelineActionConfig{
				{Name: "action-one"},
			}},
	}

	actionStore := action.NewActionStore()
	actionStore.Put("action-one", &action.MockAction{})

	jobConfigs := []config.JobConfig{
		{Name: "test job", Pipeline: "test-pipeline"},
	}

	err := worker.initPipelines(pipelineConfigs, actionStore)
	assert.NotError(t, err, "pipeline initialization failed")

	err = worker.initJobs(jobConfigs)
	assert.NotError(t, err, "job initialization failed")
}

func TestInitJobConfigFailUnknownPipeline(t *testing.T) {
	worker := &Worker{
		pipelineStore: store.New(),
	}

	pipelineConfigs := []pipeline.PipelineConfig{}

	jobConfigs := []config.JobConfig{
		{Name: "test job", Pipeline: "test-pipeline"},
	}

	err := worker.initPipelines(pipelineConfigs, nil)
	assert.NotError(t, err, "pipeline initialization failed")

	err = worker.initJobs(jobConfigs)
	assert.ErrorContains(t, err, "pipeline not found")
}

func TestSkipInitializationJobs(t *testing.T) {
	pipelineConfigs := []pipeline.PipelineConfig{
		{Name: "test-pipeline",
			Actions: []pipeline.PipelineActionConfig{
				{Name: "action-one"},
			}},
	}

	actionStore := action.NewActionStore()
	actionStore.Put("action-one", &action.MockAction{})

	jobConfigs := []config.JobConfig{
		{Name: "test job", Pipeline: "test-pipeline"},
	}

	conf := &config.Config{
		Pipelines: pipelineConfigs,
		Jobs:      jobConfigs,
	}

	worker := &Worker{
		pipelineStore: store.New(),
	}
	err := worker.init(conf, actionStore, true)
	assert.NotError(t, err, "worker creation failed")
	assert.Equal(t, 0, len(worker.jobs))
}

func testNewFromFile(t *testing.T, skipJobInitialization bool) {
	dir, err := ioutil.TempDir("/tmp", "test_worker_new_")
	assert.NotError(t, err, "creating temp dir failed")
	defer os.RemoveAll(dir)

	// ----

	clusterManager, caCert := createClusterManagerWithClusterCA(t)
	workerCert, err := clusterManager.CreateWorkerCertificate("worker-cert")
	assert.NotError(t, err, "certificate could not be created")

	caCertPath := dir + "/ca.crt"
	err = ioutil.WriteFile(caCertPath, caCert.Certificate, 0666)
	assert.NotError(t, err, "saving ca failed")

	workerCertPath := dir + "/worker.crt"
	err = ioutil.WriteFile(workerCertPath, workerCert.Certificate, 0666)
	assert.NotError(t, err, "saving certificate failed")

	workerCertKeyPath := dir + "/worker.key"
	err = ioutil.WriteFile(workerCertKeyPath, workerCert.PrivateKey, 0666)
	assert.NotError(t, err, "saving certificate key failed")

	// -----

	pipelineConfigs := []pipeline.PipelineConfig{
		{Name: "test-pipeline",
			Actions: []pipeline.PipelineActionConfig{
				{Name: "sh"},
			}},
	}

	jobConfigs := []config.JobConfig{
		{Name: "test job", Pipeline: "test-pipeline"},
	}
	conf := &config.Config{
		ServerAddr:       "server-address:server:port",
		TlsCACert:        caCertPath,
		TlsWorkerCert:    workerCertPath,
		TlsWorkerCertKey: workerCertKeyPath,
		Pipelines:        pipelineConfigs,
		Jobs:             jobConfigs,
	}

	worker, err := NewFromConfig(conf, skipJobInitialization)
	assert.NotError(t, err, "creating worker failed")
	assert.NotNil(t, worker)

	if skipJobInitialization {
		assert.Equal(t, 0, len(worker.jobs))
	}
}

func createClusterManagerWithClusterCA(t *testing.T) (manager.ClusterManager, *service.NewCertificateResponse) {
	initialClusterManager, err := manager.NewForInitialization()
	assert.NotError(t, err, "cluster manager could not be created")

	response, err := initialClusterManager.CreateClusterCACertificate("test-cluster")
	assert.NotError(t, err, "ca cert could not be created")

	clusterManager, err := manager.NewFromCA(response.Certificate, response.PrivateKey)
	assert.NotError(t, err, "ca cert could not be created")
	return clusterManager, response
}
