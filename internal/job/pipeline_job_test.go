package job

import (
	"testing"
	"time"

	"bilalekrem.com/certstore/internal/assert"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/scheduler"
	"github.com/golang/mock/gomock"
)

func TestExecute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAction := action.NewMockAction(ctrl)
	args := map[string]string{}
	args["my-arg"] = "my-value"

	mockAction.
		EXPECT().
		Run(gomock.Any(), gomock.Eq(args)).
		MinTimes(1)

	pipeline := pipeline.New("test-pipeline")
	pipeline.RegisterAction(mockAction, args)

	// ----

	mockTimeProvider := func() time.Time {
		return time.Date(2022, 01, 01, 01, 59, 58, 0, time.Local)
	}
	sched := scheduler.NewDailySchedulerWithTimeProvider(mockTimeProvider)

	pipelineJob := NewPipelineJob("test job", sched, pipeline)
	err := pipelineJob.Execute()
	assert.NotError(t, err, "while executing pipeline job")

	// wait for scheduler to run pipeline, and pipeline will run mockAction
	time.Sleep(2 * time.Second)
}
