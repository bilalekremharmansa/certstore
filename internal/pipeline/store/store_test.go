package store

import (
	"testing"

	"bilalekrem.com/certstore/internal/pipeline"
)

func TestNew(t *testing.T) {
	var firstPipeline = pipeline.New("first-pipeline")
	var secondPipeline = pipeline.New("second-pipeline")

	// ----

	store := New(firstPipeline, secondPipeline)

	retrievedFirstPipeline := store.GetPipeline("first-pipeline")
	if retrievedFirstPipeline.Name() != "first-pipeline" {
		t.Fatalf("unknown pipeline %s", retrievedFirstPipeline)
	}

	retrievedSecondPipeline := store.GetPipeline("second-pipeline")
	if retrievedSecondPipeline.Name() != "second-pipeline" {
		t.Fatalf("unknown pipeline %s", retrievedSecondPipeline)
	}

}

func TestStore(t *testing.T) {
	var pip = pipeline.New("first-pipeline")

	// ----

	store := New()
	store.StorePipeline(pip)

	retrievedFirstPipeline := store.GetPipeline("first-pipeline")
	if retrievedFirstPipeline.Name() != "first-pipeline" {
		t.Fatalf("unknown pipeline %s", retrievedFirstPipeline)
	}

}
