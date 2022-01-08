package pipeline

import (
	"testing"
)

func TestParseFromByte(t *testing.T) {
	pipelineYaml := `name: my-pipeline
actions:
  - name: shell-cmd
    args:
      command: "echo hello"
  - name: mock-action`

	config, err := ParsePipelineConfig(pipelineYaml)
	if err != nil {
		t.Fatalf("Error occurred while parsing pipeline config, %v\n", err)
	}

	if config.Name != "my-pipeline" {
		t.Fatalf("Pipeline name is not correct, expected: my-pipeline found: %s\n", config.Name)
	}

	actions := config.Actions
	if len(actions) != 2 {
		t.Fatalf("pipeline action length is not correct, expected: 2, found: %d\n", len(actions))
	}

	// ----

	firstAction := actions[0]
	if firstAction.Name != "shell-cmd" {
		t.Fatalf("first action name is not correct, expected: shell-cmd, found: %s\n", firstAction.Name)
	}

	if len(firstAction.Args) != 1 {
		t.Fatalf("first action arg length is not correct, expected: 1, found: %d\n", len(firstAction.Args))
	}

	if firstAction.Args["command"] != "echo hello" {
		t.Fatalf("first action argument is not correct, expected: 'command=echo hello', found: %s\n", firstAction.Args)
	}

	// ----

	secondAction := actions[1]
	if secondAction.Name != "mock-action" {
		t.Fatalf("second action name is not correct, expected: mock-action, found: %ss\n", secondAction.Name)
	}

	// should be a nil map
	if secondAction.Args != nil {
		t.Fatalf("second action argument is not correct, expected: nil, found: %s\n", secondAction.Args)
	}

}
