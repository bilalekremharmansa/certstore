package pipeline

import (
	"errors"
	"fmt"

	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/context"
	"bilalekrem.com/certstore/internal/pipeline/store"
)

const (
	ARGS_PIPELINE_NAME string = "pipeline-name"
)

type pipelineAction struct {
	store *store.PipelineStore
}

func NewPipelineAction(store *store.PipelineStore) pipelineAction {
	return pipelineAction{store: store}
}

func (a pipelineAction) Run(ctx *context.Context, args map[string]string) error {
	err := validate(args)
	if err != nil {
		logging.GetLogger().Errorf("validation args failed, %v", err)
		return err
	}

	// ----

	pipelineName := args[ARGS_PIPELINE_NAME]

	pipelineStore := a.store
	pipeline := pipelineStore.GetPipeline(pipelineName)
	if pipeline == nil {
		logging.GetLogger().Errorf("pipeline not in pipeline store: [%s]", pipelineName)
		return errors.New("pipeline not found")
	}

	// ----

	err = pipeline.Run()
	if err != nil {
		logging.GetLogger().Errorf("Running pipeline failed in pipeline action, %v", err)
		return err
	}

	return nil
}

func validate(args map[string]string) error {
	_, exists := args[ARGS_PIPELINE_NAME]
	if !exists {
		return errors.New(fmt.Sprintf("required argument: %s", ARGS_PIPELINE_NAME))
	}

	return nil
}
