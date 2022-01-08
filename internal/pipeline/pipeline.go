package pipeline

import (
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline/action"
)

type Pipeline interface {
	Name() string
	Run() error
}

// ----

type pipelineImpl struct {
	name    string
	actions []*ActionWithArgs
}

func New(name string) *pipelineImpl {
	return &pipelineImpl{name: name}
}

func NewFromConfig(config *PipelineConfig, actionStore *action.ActionStore) (*pipelineImpl, error) {
	pipeline := New(config.Name)

	for _, configAction := range config.Actions {
		_action, err := actionStore.Get(configAction.Name)

		if err != nil {
			return nil, err
		}

		args := configAction.Args
		logging.GetLogger().Infof("[%s] args [%s]", configAction.Name, args)
		pipeline.RegisterAction(*_action, args)
	}

	return pipeline, nil
}

func (p *pipelineImpl) RegisterAction(action action.Action, args map[string]string) {
	logging.GetLogger().Infof("Adding a new action to pipeline [%s]", p.Name())
	actionWithArgs := ActionWithArgs{action: action, args: args}
	p.actions = append(p.actions, &actionWithArgs)
}

func (p *pipelineImpl) Name() string {
	return p.name
}

func (p *pipelineImpl) Run() error {
	logging.GetLogger().Infof("Starting to run pipeline [%s]", p.Name())
	for _, actionWithArgs := range p.actions {
		action := actionWithArgs.action
		args := actionWithArgs.args
		err := action.Run(args)
		if err != nil {
			logging.GetLogger().Debugf("error occurred while running an action in pipeline %s, err: %v", p.Name(), err)
			return err
		}
	}
	return nil
}
