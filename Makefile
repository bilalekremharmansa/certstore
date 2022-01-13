.PHONY: generate-proto

generate-proto:
	@PATH="$(PATH):$(go env GOPATH)/bin" $(shell protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/grpc/proto/*.proto)
	@PATH="$(PATH):$(go env GOPATH)/bin" $(shell protoc --go_out=internal/certstore/grpc/gen --go_opt=paths=source_relative --go-grpc_out=internal/certstore/grpc/gen --go-grpc_opt=paths=source_relative --proto_path=internal/certstore/grpc/proto internal/certstore/grpc/proto/*proto)


generate-mock:
	mockgen -source=internal/pipeline/action/action.go -destination internal/pipeline/action/mock_action.go -package=action
	mockgen -source=internal/certificate/service/service.go -destination internal/certificate/service/mock_service.go -package=service
	mockgen -source=internal/certstore/certstore.go -destination internal/certstore/mock_certstore.go -package=certstore

test:
	go test ./...

