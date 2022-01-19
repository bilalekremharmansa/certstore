package worker

import (
	"testing"

	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
)

func TestInitSuccess(t *testing.T) {
	worker := &Worker{
		pipelines: make(map[string]pipeline.Pipeline),
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
	err := worker.init(pipelineConfigs, actionStore)
	if err != nil {
		t.Fatalf("initialization failed, %v", err)
	}
}

func TestInitFail(t *testing.T) {
	worker := &Worker{
		pipelines: make(map[string]pipeline.Pipeline),
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
	err := worker.init(pipelineConfigs, actionStore)
	if err == nil {
		t.Fatalf("initialization should've been failed -- action-two is missing in store, %v", err)
	}
}
