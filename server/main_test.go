package main

import (
	"fmt"
	"net"
	"os"
	"testing"
)

const (
	envTestPort = "TEST_PORT"
)

func TestRunMain(t *testing.T) {
	port := os.Getenv(envTestPort)
	if port == "" {
		t.Fatal("must specify port number")
	}

	go main()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		t.Fatal("failed to listen to the termination request")
	}

	conn, err := lis.Accept()
	if err != nil {
		t.Fatal("failed to accept the termination request")
	}

	err = conn.Close()
	if err != nil {
		t.Fatal("failed to close termination connection")
	}
}
