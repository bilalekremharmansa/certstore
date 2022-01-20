package worker

import (
	"errors"
	"fmt"
	"io/ioutil"

	certificate_service "bilalekrem.com/certstore/internal/certstore/grpc/gen"
	"bilalekrem.com/certstore/internal/cluster/worker/config"
	"bilalekrem.com/certstore/internal/logging"
	"bilalekrem.com/certstore/internal/pipeline"
	"bilalekrem.com/certstore/internal/pipeline/action"
	"bilalekrem.com/certstore/internal/pipeline/action/issuecertificate"
	"google.golang.org/grpc"
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

	logging.GetLogger().Info("creating certificate service client for action store")
	certificateServiceClient, err := getCertificateServiceClient(&conf.Cluster)
	if err != nil {
		logging.GetLogger().Errorf("creating cert service client faild, %v", err)
		return nil, err
	}

	actionStore := getActionStore(certificateServiceClient)
	worker.init(conf.Pipelines, actionStore)

	// ----

	return worker, nil
}

func (w *Worker) RunPipeline(pipelineName string) error {
	logging.GetLogger().Errorf("Running pipeline, %s", pipelineName)
	pip, exists := w.pipelines[pipelineName]
	if !exists {
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

func (w *Worker) init(pipelineConfigs []pipeline.PipelineConfig, actionStore *action.ActionStore) error {
	logging.GetLogger().Info("Initializing worker with pipeline configs..")
	for _, pipelineConfig := range pipelineConfigs {
		pip, err := pipeline.NewFromConfig(&pipelineConfig, actionStore)
		if err != nil {
			logging.GetLogger().Errorf("creating pipeline failed, %v", err)
			return errors.New(fmt.Sprintf("creating pipeline failed, %v", err))
		}

		logging.GetLogger().Info("pipeline is created: [%s]", pip.Name())
		w.pipelines[pip.Name()] = pip
	}

	return nil
}

func getActionStore(client *certificate_service.CertificateServiceClient) *action.ActionStore {
	store := action.NewActionStore()

	store.Put("issue-certificate", issuecertificate.NewIssueCertificateAction(*client))

	return store
}
