package pipeline

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
)

func TestParseFromByte(t *testing.T) {
	pipelineYaml := `name: my-pipeline
actions:
  - name: shell-cmd
    args:
      command: "echo hello"
  - name: mock-action`

	config, err := ParsePipelineConfig(pipelineYaml)
	assert.NotError(t, err, "parsing pipeline config")

	assert.Equal(t, "my-pipeline", config.Name)

	actions := config.Actions
	assert.Equal(t, 2, len(actions))

	// ----

	firstAction := actions[0]
	assert.Equal(t, "shell-cmd", firstAction.Name)
	assert.Equal(t, 1, len(firstAction.Args))
	assert.Equal(t, "echo hello", firstAction.Args["command"])

	// ----

	secondAction := actions[1]
	assert.Equal(t, "mock-action", secondAction.Name)

	// should be a nil map
	assert.Nil(t, secondAction.Args)
}
