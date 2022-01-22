package store

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/pipeline"
)

func TestNew(t *testing.T) {
	var firstPipeline = pipeline.New("first-pipeline")
	var secondPipeline = pipeline.New("second-pipeline")

	// ----

	store := New(firstPipeline, secondPipeline)

	retrievedFirstPipeline := store.GetPipeline("first-pipeline")
	assert.Equal(t, "first-pipeline", retrievedFirstPipeline.Name())

	retrievedSecondPipeline := store.GetPipeline("second-pipeline")
	assert.Equal(t, "second-pipeline", retrievedSecondPipeline.Name())
}

func TestStore(t *testing.T) {
	var pip = pipeline.New("first-pipeline")

	// ----

	store := New()
	store.StorePipeline(pip)

	retrievedFirstPipeline := store.GetPipeline("first-pipeline")
	assert.Equal(t, "first-pipeline", retrievedFirstPipeline.Name())

}
