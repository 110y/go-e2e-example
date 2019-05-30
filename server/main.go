package main

import (
	"os"

	"fmt"
	"net"

	"os/signal"
	"syscall"

	"github.com/110y/go-e2e-example/server/pb"
	"google.golang.org/grpc"
)

const (
	envPort                = "PORT"
	statusPortNotSpecified = 1
	statusFailedToListen   = 2
	statusFailedToServe    = 3
)

func main() {
	port := os.Getenv(envPort)
	if port == "" {
		fmt.Fprintf(os.Stderr, "must specify %s envionment variable\n", envPort)
		os.Exit(statusPortNotSpecified)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen: %+v", err)
		os.Exit(statusFailedToListen)
	}

	gs := grpc.NewServer()
	pb.RegisterServerServer(gs, &server{})

	go func() {
		if err := gs.Serve(lis); err != nil {
			fmt.Fprintf(os.Stderr, "failed to serve: %+v", err)
			os.Exit(statusFailedToServe)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	<-sigChan
	gs.GracefulStop()
}
