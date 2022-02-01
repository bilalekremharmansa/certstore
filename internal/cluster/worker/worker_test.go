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
	err := worker.init(pipelineConfigs, actionStore, []config.JobConfig{})
	assert.NotError(t, err, "initialization failed")
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
	err := worker.init(pipelineConfigs, actionStore, []config.JobConfig{})
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

	err := worker.init(pipelineConfigs, actionStore, jobConfigs)
	assert.NotError(t, err, "initialization failed")
}

func TestInitJobConfigFailUnknownPipeline(t *testing.T) {
	worker := &Worker{
		pipelineStore: store.New(),
	}

	pipelineConfigs := []pipeline.PipelineConfig{}

	jobConfigs := []config.JobConfig{
		{Name: "test job", Pipeline: "test-pipeline"},
	}

	err := worker.init(pipelineConfigs, nil, jobConfigs)
	assert.ErrorContains(t, err, "pipeline not found")
}
