package agent

import (
	"io/ioutil"
	"os"
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/certificate/service"
	"bilalekrem.com/certstore/internal/cluster/manager"
	"bilalekrem.com/certstore/internal/cluster/agent/config"
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
	agent := &Agent{
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
	err := agent.initPipelines(pipelineConfigs, actionStore)
	assert.NotError(t, err, "pipeline initialization failed")
}

func TestInitPipelineConfigFail(t *testing.T) {
	agent := &Agent{
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
	err := agent.initPipelines(pipelineConfigs, actionStore)
	assert.Error(t, err, "should've been failed, action-two is missing in store")
}

func TestInitJobConfigSuccess(t *testing.T) {
	agent := &Agent{
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

	err := agent.initPipelines(pipelineConfigs, actionStore)
	assert.NotError(t, err, "pipeline initialization failed")

	err = agent.initJobs(jobConfigs)
	assert.NotError(t, err, "job initialization failed")
}

func TestInitJobConfigFailUnknownPipeline(t *testing.T) {
	agent := &Agent{
		pipelineStore: store.New(),
	}

	pipelineConfigs := []pipeline.PipelineConfig{}

	jobConfigs := []config.JobConfig{
		{Name: "test job", Pipeline: "test-pipeline"},
	}

	err := agent.initPipelines(pipelineConfigs, nil)
	assert.NotError(t, err, "pipeline initialization failed")

	err = agent.initJobs(jobConfigs)
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

	agent := &Agent{
		pipelineStore: store.New(),
	}
	err := agent.init(conf, actionStore, true)
	assert.NotError(t, err, "agent creation failed")
	assert.Equal(t, 0, len(agent.jobs))
}

func testNewFromFile(t *testing.T, skipJobInitialization bool) {
	dir, err := ioutil.TempDir("/tmp", "test_agent_new_")
	assert.NotError(t, err, "creating temp dir failed")
	defer os.RemoveAll(dir)

	// ----

	clusterManager, caCert := createClusterManagerWithClusterCA(t)
	agentCert, err := clusterManager.CreateAgentCertificate("agent-cert")
	assert.NotError(t, err, "certificate could not be created")

	caCertPath := dir + "/ca.crt"
	err = ioutil.WriteFile(caCertPath, caCert.Certificate, 0666)
	assert.NotError(t, err, "saving ca failed")

	agentCertPath := dir + "/agent.crt"
	err = ioutil.WriteFile(agentCertPath, agentCert.Certificate, 0666)
	assert.NotError(t, err, "saving certificate failed")

	agentCertKeyPath := dir + "/agent.key"
	err = ioutil.WriteFile(agentCertKeyPath, agentCert.PrivateKey, 0666)
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
		TlsAgentCert:    agentCertPath,
		TlsAgentCertKey: agentCertKeyPath,
		Pipelines:        pipelineConfigs,
		Jobs:             jobConfigs,
	}

	agent, err := NewFromConfig(conf, skipJobInitialization)
	assert.NotError(t, err, "creating agent failed")
	assert.NotNil(t, agent)

	if skipJobInitialization {
		assert.Equal(t, 0, len(agent.jobs))
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
