package worker

import (
	"errors"
	"fmt"
	"io/ioutil"

	"bilalekrem.com/certstore/internal/cluster/worker/config"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
)

type Worker struct {
	pipelines map[string]pipeline.Pipeline
}

func NewFromFile(path string) (*Worker, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err := config.Parse(string(bytes))
	if err != nil {
		return nil, err
	}

	return NewFromConfig(config)
}

func NewFromConfig(conf *config.Config) (*Worker, error) {
	worker := &Worker{
		pipelines: make(map[string]pipeline.Pipeline),
	}

	// ----

	actionStore := getActionStore()
	worker.init(conf.Pipelines, actionStore)

	// ----

	return worker, nil
}

func (w *Worker) init(pipelineConfigs []pipeline.PipelineConfig, actionStore *action.ActionStore) error {
	for _, pipelineConfig := range pipelineConfigs {
		pip, err := pipeline.NewFromConfig(&pipelineConfig, actionStore)
		if err != nil {
			return errors.New(fmt.Sprintf("creating pipeline failed, %v", err))
		}

		w.pipelines[pipelineConfig.Name] = pip
	}

	return nil
}

func getActionStore() *action.ActionStore {
	store := action.NewActionStore()

	return store
}
