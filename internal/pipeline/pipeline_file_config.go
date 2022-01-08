package pipeline

import (
	"gopkg.in/yaml.v3"
)

/*
sample pipeline config:

name: my-pipeline
actions:
  - name: shell-cmd
    args:
	  command: "echo hello world"
  - name: mock-action
*/

type PipelineConfig struct {
	Name    string                 `yaml:"name"`
	Actions []PipelineActionConfig `yaml:"actions"`
}

type PipelineActionConfig struct {
	Name string            `yaml:"name"`
	Args map[string]string `yaml:"args"`
}

func ParsePipelineConfig(configYaml string) (*PipelineConfig, error) {
	config := &PipelineConfig{}
	configBytes := []byte(configYaml)
	err := yaml.Unmarshal(configBytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
