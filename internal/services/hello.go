package services

import (
	"context"
	"log"

	"github.com/daheige/gorpc/api/clients/go/pb"
)

// GreeterService rpc service entry
type GreeterService struct{}

// SayHello say hello
func (s *GreeterService) SayHello(ctx context.Context, in *pb.HelloReq) (*pb.HelloReply, error) {
	log.Println("req data: ", in)

	// mock business logic
	if in.Id == 1 && in.Name == "" {
		in.Name = "micro"
	}

	return &pb.HelloReply{
		Name:    "hello," + in.Name,
		Message: "call ok",
	}, nil
}
