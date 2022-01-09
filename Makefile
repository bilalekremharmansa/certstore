.PHONY: generate-proto

generate-proto:
	@PATH="$(PATH):$(go env GOPATH)/bin" $(shell protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/grpc/proto/*.proto)


generate-mock:
	mockgen -source=internal/pipeline/action/action.go -destination internal/pipeline/action/mock_action.go -package=action

test:
	go test ./...

