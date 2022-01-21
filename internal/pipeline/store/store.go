package store

import "bilalekrem.com/certstore/internal/pipeline"

type PipelineStore struct {
	pipelines map[string]pipeline.Pipeline
}

func New(pipelines ...pipeline.Pipeline) *PipelineStore {
	store := &PipelineStore{make(map[string]pipeline.Pipeline)}

	for _, pip := range pipelines {
		store.StorePipeline(pip)
	}

	return store
}

func (p *PipelineStore) GetPipeline(key string) pipeline.Pipeline {
	value, exist := p.pipelines[key]
	if !exist {
		return nil
	}

	return value
}

func (p *PipelineStore) StorePipeline(pip pipeline.Pipeline) {
	p.pipelines[pip.Name()] = pip
}
