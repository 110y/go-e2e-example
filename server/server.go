package main

import (
	"context"

	"github.com/110y/go-e2e-example/server/pb"
)

type server struct{}

func (s *server) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Echo: req.Echo,
	}, nil
}
