package worker

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/cluster/worker/config"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/store"
)

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
	assert.Nil(t, worker.jobs)
}
