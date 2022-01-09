package pipeline

import (
	"reflect"
	"testing"

	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/context"
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
		Run(gomock.Any(), gomock.Any()).
		MinTimes(1)

	pipeline.RegisterAction(mockAction, nil)

	// ----

	err := pipeline.Run()
	if err != nil {
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
		Run(gomock.Any(), gomock.Eq(args)).
		MinTimes(1)
	pipeline.RegisterAction(mockAction, args)

	// ----

	err := pipeline.Run()
	if err != nil {
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
		Run(gomock.Any(), gomock.Any()).
		MinTimes(1)

	mockAction2 := action.NewMockAction(ctrl)
	mockAction2.
		EXPECT().
		Run(gomock.Any(), gomock.Any()).
		MinTimes(1)

	pipeline.RegisterAction(mockAction, nil)
	pipeline.RegisterAction(mockAction2, nil)

	// ----

	err := pipeline.Run()
	if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}

}

func TestNewPipelineFromConfig(t *testing.T) {
	actionsConfig := []PipelineActionConfig{
		{Name: "my-action", Args: nil},
	}
	pipelineConfig := &PipelineConfig{Name: "my-pipeline", Actions: actionsConfig}

	// ----

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := action.NewMockAction(ctrl)
	mockAction.
		EXPECT().
		Run(gomock.Any(), gomock.Any()).
		MinTimes(1)

	actionStore := action.NewActionStore()
	actionStore.Put("my-action", mockAction)

	// ----

	pipeline, err := NewFromConfig(pipelineConfig, actionStore)
	if err != nil {
		t.Fatalf("error occurred while initation pipeline from pipeline config, %v", err)
	}

	// ----

	err = pipeline.Run()
	if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}

}

func TestNewPipelineFromConfigMissingAction(t *testing.T) {
	actionsConfig := []PipelineActionConfig{
		{Name: "my-action", Args: nil},
	}
	pipelineConfig := &PipelineConfig{Name: "my-pipeline", Actions: actionsConfig}

	actionStore := action.NewActionStore()
	_, err := NewFromConfig(pipelineConfig, actionStore)
	if err == nil {
		t.Fatalf("error is expected beceause of missing action in action store, but not found")
	}
}

func TestNewPipelineFromYamlConfig(t *testing.T) {
	pipelineYaml := `name: my-pipeline
actions:
  - name: shell-cmd
    args:
      command: "echo hello"
  - name: test-action`

	pipelineConfig, err := ParsePipelineConfig(pipelineYaml)
	if err != nil {
		t.Fatalf("error occurred while parsing pipeline pipeline config, %v", err)
	}

	// -----

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shellArgs := map[string]string{}
	shellArgs["command"] = "echo hello"

	shellCommandAction := action.NewMockAction(ctrl)
	shellCommandAction.
		EXPECT().
		Run(gomock.Any(), gomock.Eq(shellArgs)).
		MinTimes(0)

	var nilMap map[string]string
	testAction := action.NewMockAction(ctrl)
	testAction.
		EXPECT().
		Run(gomock.Any(), gomock.Eq(nilMap)).
		MinTimes(0)

	actionStore := action.NewActionStore()
	actionStore.Put("shell-cmd", shellCommandAction)
	actionStore.Put("test-action", testAction)

	// -----

	pipeline, err := NewFromConfig(pipelineConfig, actionStore)
	if err != nil {
		t.Fatalf("error occurred while initation pipeline from pipeline config, %v", err)
	}

	// ----

	err = pipeline.Run()
	if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}
}

func TestPipelineContextStoreAndGetCustomValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// ----

	var FIRST_ACTION_STRING_KEY context.Key = "string key"
	FIRST_ACTION_STRING_VALUE := "hello world"

	var FIRST_ACTION_NUMBER_KEY context.Key = "number key"
	FIRST_ACTION_NUMBER_VALUE := 100

	first := action.NewMockAction(ctrl)
	first.
		EXPECT().
		Run(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx *context.Context, args map[string]string) error {
			ctx.StoreValue(FIRST_ACTION_STRING_KEY, FIRST_ACTION_STRING_VALUE)
			ctx.StoreValue(FIRST_ACTION_NUMBER_KEY, FIRST_ACTION_NUMBER_VALUE)

			return nil
		}).
		AnyTimes()

	second := action.NewMockAction(ctrl)
	second.
		EXPECT().
		Run(gomock.Any(), gomock.Any()).
		Do(func(ctx *context.Context, args map[string]string) error {
			stringValue := ctx.GetValue(FIRST_ACTION_STRING_KEY)

			if reflect.TypeOf(stringValue).Kind() != reflect.String {
				t.Fatalf("Unexpected type, expected: string, found:%T", reflect.TypeOf(stringValue))
			}

			if stringValue != FIRST_ACTION_STRING_VALUE {
				t.Fatalf("Unexpected value, expected: %s, found:%v", FIRST_ACTION_STRING_VALUE, stringValue)
			}

			// -----

			numberValue := ctx.GetValue(FIRST_ACTION_NUMBER_KEY)

			if reflect.TypeOf(numberValue).Kind() != reflect.Int {
				t.Fatalf("Unexpected type, expected: int, found:%T", reflect.TypeOf(numberValue))
			}

			if numberValue != FIRST_ACTION_NUMBER_VALUE {
				t.Fatalf("Unexpected value, expected: %d, found:%v", FIRST_ACTION_NUMBER_VALUE, numberValue)
			}

			return nil
		}).
		AnyTimes()

	pipeline := New("test-pipeline")
	pipeline.RegisterAction(first, nil)
	pipeline.RegisterAction(second, nil)

	err := pipeline.Run()
	if err != nil {
		t.Fatalf("Error occurred while running pipeline, %v", err)
	}
}
