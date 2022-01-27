package job

import (
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/scheduler"
)

type pipelineJob struct {
	name      string
	scheduler scheduler.Scheduler
	pipeline  pipeline.Pipeline
}

func NewPipelineJob(name string, sched scheduler.Scheduler, pip pipeline.Pipeline) *pipelineJob {
	return &pipelineJob{
		name:      name,
		scheduler: sched,
		pipeline:  pip,
	}
}

func (j *pipelineJob) Execute() error {
	logging.GetLogger().Infof("Scheduling job [%s]", j.name)

	j.scheduler.Schedule(func() {
		err := j.pipeline.Run()
		logging.GetLogger().Errorf("Running job [%s] pipeline failed, %v", j.name, err)
	})

	return nil
}
