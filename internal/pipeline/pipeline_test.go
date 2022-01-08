package pipeline

import (
	"testing"

	"bilalekrem.com/certstore/internal/pipeline/action"
	"github.com/golang/mock/gomock"
)

func TestPipelineName(t *testing.T) {
	pipelineName := "test-pipeline"
	pipeline := New(pipelineName)

	if pipeline.Name() != pipelineName {
		t.Fatalf("pipeline name is not correct, expected: [%s], actual: [%s]", pipelineName, pipeline.Name())
	}
}

func TestPipelineRunAction(t *testing.T) {
	pipelineName := "test-pipeline"
	pipeline := New(pipelineName)

	// ----

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := action.NewMockAction(ctrl)
	mockAction.
		EXPECT().
		Run(gomock.Any()).
		MinTimes(1)

	pipeline.RegisterAction(mockAction, nil)

	// ----

	err := pipeline.Run(); if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}

}

func TestPipelineRunActionWithConfig(t *testing.T) {
	pipelineName := "test-pipeline"
	pipeline := New(pipelineName)

	// ----

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := action.NewMockAction(ctrl)
	args := map[string]string{}
	args["my-arg"] = "my-value"

	mockAction.
		EXPECT().
		Run(gomock.Eq(args)).
		MinTimes(1)
	pipeline.RegisterAction(mockAction, args)

	// ----

	err := pipeline.Run(); if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}

}

func TestPipelineRunMultipleAction(t *testing.T) {
	pipelineName := "test-pipeline"
	pipeline := New(pipelineName)

	// ----

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := action.NewMockAction(ctrl)
	mockAction.
		EXPECT().
		Run(gomock.Any()).
		MinTimes(1)

	mockAction2 := action.NewMockAction(ctrl)
	mockAction2.
		EXPECT().
		Run(gomock.Any()).
		MinTimes(1)

	pipeline.RegisterAction(mockAction, nil)
	pipeline.RegisterAction(mockAction2, nil)

	// ----

	err := pipeline.Run(); if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}

}
