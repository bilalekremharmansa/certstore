package worker

import (
	"errors"
	"fmt"
	"io/ioutil"

	certificate_service "bilalekrem.com/certstore/internal/certstore/grpc/gen"
	"bilalekrem.com/certstore/internal/cluster/worker/config"
	"bilalekrem.com/certstore/internal/job"
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/action/issuecertificate"
	pipeline_action "bilalekrem.com/certstore/internal/pipeline/action/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action/savecertificate"
	"bilalekrem.com/certstore/internal/pipeline/action/shouldrenewcertificate"
	"bilalekrem.com/certstore/internal/pipeline/store"
	"bilalekrem.com/certstore/internal/scheduler"
	"google.golang.org/grpc"
)

type Worker struct {
	pipelineStore *store.PipelineStore
	jobs          []job.Job
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
		pipelineStore: store.New(),
		jobs:          []job.Job{},
	}

	// ----

	logging.GetLogger().Info("creating certificate service client for action store")
	certificateServiceClient, err := getCertificateServiceClient(&conf.Cluster)
	if err != nil {
		logging.GetLogger().Errorf("creating cert service client faild, %v", err)
		return nil, err
	}

	actionStore := getActionStore(certificateServiceClient, worker.pipelineStore)
	worker.init(conf, actionStore)

	// ----

	return worker, nil
}

func (w *Worker) RunPipeline(pipelineName string) error {
	logging.GetLogger().Errorf("Running pipeline, %s", pipelineName)
	pip := w.pipelineStore.GetPipeline(pipelineName)
	if pip == nil {
		logging.GetLogger().Errorf("pipeline not found, %s", pipelineName)
		return errors.New(fmt.Sprintf("pipeline not found, %s", pipelineName))
	}

	return pip.Run()
}

// ----

func getCertificateServiceClient(conf *config.ClusterConfig) (*certificate_service.CertificateServiceClient, error) {
	serverAddress := conf.ServerAddr

	// ---

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		return nil, err
	}
	// todo bilal defer conn.Close()

	client := certificate_service.NewCertificateServiceClient(conn)
	return &client, nil
}

func (w *Worker) init(conf *config.Config, actionStore *action.ActionStore) error {
	err := w.initPipelines(conf.Pipelines, actionStore)
	if err != nil {
		return err
	}

	err = w.initJobs(conf.Jobs)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) initPipelines(pipelineConfigs []pipeline.PipelineConfig, actionStore *action.ActionStore) error {
	logging.GetLogger().Info("Initializing worker with pipeline configs..")
	for _, pipelineConfig := range pipelineConfigs {
		pip, err := pipeline.NewFromConfig(&pipelineConfig, actionStore)
		if err != nil {
			logging.GetLogger().Errorf("creating pipeline failed, %v", err)
			return errors.New(fmt.Sprintf("creating pipeline failed, %v", err))
		}

		logging.GetLogger().Infof("pipeline is created: [%s]", pip.Name())
		w.pipelineStore.StorePipeline(pip)
	}

	return nil
}

func (w *Worker) initJobs(jobConfigs []config.JobConfig) error {
	logging.GetLogger().Info("Initializing worker with job configs..")
	for _, jobConfig := range jobConfigs {
		dailyScheduler := scheduler.NewDailyScheduler()

		pip := w.pipelineStore.GetPipeline(jobConfig.Pipeline)
		if pip == nil {
			logging.GetLogger().Errorf("pipeline not found, %s", jobConfig.Pipeline)
			return errors.New(fmt.Sprintf("pipeline not found, %s", jobConfig.Pipeline))
		}
		pipelineJob := job.NewPipelineJob(jobConfig.Name, dailyScheduler, pip)
		err := pipelineJob.Execute()
		if err != nil {
			logging.GetLogger().Errorf("Running job failed, %s", jobConfig.Pipeline)
			return err
		}
		w.jobs = append(w.jobs, pipelineJob)
	}

	return nil
}

func getActionStore(client *certificate_service.CertificateServiceClient, pipelineStore *store.PipelineStore) *action.ActionStore {
	store := action.NewActionStore()

	store.Put("issue-certificate", issuecertificate.NewIssueCertificateAction(*client))
	store.Put("save-certificate", savecertificate.NewSaveCertificateAction())
	store.Put("run-pipeline", pipeline_action.NewPipelineAction(pipelineStore))
	store.Put("should-renew-certificate", shouldrenewcertificate.NewShouldRenewCertificateAction())

	return store
}
