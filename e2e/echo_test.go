package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/110y/go-e2e-example/server/pb"
)

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
