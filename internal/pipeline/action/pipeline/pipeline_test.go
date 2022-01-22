package pipeline

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/context"
	"bilalekrem.com/certstore/internal/pipeline/store"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// ----

	mockAction := action.NewMockAction(ctrl)
	mockAction.
		EXPECT().
		Run(gomock.Any(), gomock.Any()).
		MinTimes(1)

	pip := pipeline.New("my-pipeline")
	pip.RegisterAction(mockAction, nil)

	pipelineStore := store.New(pip)

	// ----

	action := NewPipelineAction(pipelineStore)

	args := make(map[string]string)
	args[ARGS_PIPELINE_NAME] = "my-pipeline"

	err := action.Run(context.New(), args)
	assert.NotError(t, err, "running action")
}

func TestRequiredPipelineName(t *testing.T) {
	args := make(map[string]string)

	err := NewPipelineAction(nil).Run(nil, args)
	assert.ErrorContains(t, err, "required argument")
}

func TestMissingPipeline(t *testing.T) {
	pipelineStore := store.New()

	args := make(map[string]string)
	args[ARGS_PIPELINE_NAME] = "my-pipeline"

	err := NewPipelineAction(pipelineStore).Run(nil, args)
	assert.ErrorContains(t, err, "pipeline not found")
}
