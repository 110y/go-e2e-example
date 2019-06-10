package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/110y/go-e2e-example/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcClient pb.ServerClient

const (
	envPort                = "PORT"
	statusPortNotSpecified = 1
	statusFailedToListen   = 2
	statusFailedToServe    = 3
)

func TestMain(m *testing.M) {
	os.Exit(func() (status int) {
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
		pb.RegisterServerServer(gs, &Server{})
		reflection.Register(gs)

		go func() {
			if err := gs.Serve(lis); err != nil {
				fmt.Fprintf(os.Stderr, "failed to serve: %+v", err)
				os.Exit(statusFailedToServe)
			}
		}()

		conn, err := grpc.DialContext(
			context.Background(),
			fmt.Sprintf(":%s", port),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithTimeout(5*time.Second),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to the test server: %+v\n", err)
			return 1
		}
		grpcClient = pb.NewServerClient(conn)

		defer func() {
			err := conn.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to close the connection to the test server: %+v\n", err)
				status = 1
			}
		}()

		return m.Run()
	}())
}

func TestEcho(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.EchoRequest{
		Echo: "foo",
	}

	res, err := grpcClient.Echo(ctx, req)

	if err != nil {
		t.Fatalf("failed to request to the server: %s", err)
	}

	if res.Echo != req.Echo {
		t.Errorf("want %s, but got %s", req.Echo, res.Echo)
	}
}
