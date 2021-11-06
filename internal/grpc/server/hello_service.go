package server

import (
	"context"

	pb "bilalekrem.com/certstore/internal/grpc/proto"
)

type helloService struct {
	pb.UnimplementedHelloServiceServer
}

func NewHelloService() *helloService {
	return &helloService{}
}

func (s *helloService) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello " + req.Name,
	}, nil
}
