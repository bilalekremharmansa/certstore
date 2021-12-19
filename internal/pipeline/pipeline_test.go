package pipeline

import (
	"testing"

	"bilalekrem.com/certstore/internal/pipeline/action"
)

func TestPipelineName(t *testing.T) {
	pipelineName := "test-pipeline"
	pipeline := New(pipelineName)

	if pipeline.Name() != pipelineName {
		t.Fatalf("pipeline name is not correct, expected: [%s], actual: [%s]", pipelineName, pipeline.Name())
	}
}

func TestPipelineRunMultipleAction(t *testing.T) {
	pipelineName := "test-pipeline"
	pipeline := New(pipelineName)

	// ----

	mockAction := &action.MockAction{}
	pipeline.RegisterAction(mockAction, nil)

	mockAction2 := &action.MockAction{}
	pipeline.RegisterAction(mockAction2, nil)

	// ----

	err := pipeline.Run(); if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}

	// ----

	if !mockAction.Executed {
		t.Fatalf("action should've been executed, but did not")
	}
	if !mockAction2.Executed {
		t.Fatalf("action should've been executed, but did not")
	}
}