include .env
export

.PHONY: generate-proto build

define build_go_file
	GOOS=windows GOARCH=amd64 go build -o build/certstore.exe cmd/main.go
	GOOS=linux go build -o build/cerstore cmd/main.go
endef

generate-proto:
	@PATH="$(PATH):$(go env GOPATH)/bin" $(shell protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/grpc/proto/*.proto)
	@PATH="$(PATH):$(go env GOPATH)/bin" $(shell protoc --go_out=internal/certstore/grpc/gen --go_opt=paths=source_relative --go-grpc_out=internal/certstore/grpc/gen --go-grpc_opt=paths=source_relative --proto_path=internal/certstore/grpc/proto internal/certstore/grpc/proto/*proto)


generate-mock:
	mockgen -source=internal/pipeline/action/action.go -destination internal/pipeline/action/mock_action.go -package=action
	mockgen -source=internal/certificate/service/service.go -destination internal/certificate/service/mock_service.go -package=service
	mockgen -source=internal/certstore/certstore.go -destination internal/certstore/mock_certstore.go -package=certstore
	mockgen -source=internal/certstore/grpc/gen/certificate_service_grpc.pb.go -destination internal/certstore/grpc/gen/mock_certificate_service.pb.go -package=gen
	mockgen -source=internal/lego/adapter.go -destination internal/lego/mock_adapter.go -package=lego

test:
	go test ./...

build:
	$(call build_go_file)

run-server:
	go run cmd/main.go cluster server start --config $(SERVER_CONFIG_PATH)

run-pipeline:
	go run cmd/main.go cluster worker runPipeline --config $(WORKER_CONFIG_PATH) --pipeline $(WORKER_PIPELINE_TO_RUN)

run-worker:
	go run cmd/main.go cluster worker start --config $(WORKER_CONFIG_PATH)
