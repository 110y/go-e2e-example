package e2e

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/110y/go-e2e-example/server/pb"
	"google.golang.org/grpc"
)

var grpcClient pb.ServerClient

const (
	envPort     = "PORT"
	envTestPort = "TEST_PORT"
)

func TestMain(m *testing.M) {
	os.Exit(func() (status int) {
		port := os.Getenv(envPort)
		if port == "" {
			fmt.Fprintf(os.Stderr, "must specify %s envionment variable\n", envPort)
			return 1
		}

		testPort := os.Getenv(envTestPort)
		if testPort == "" {
			fmt.Fprintf(os.Stderr, "must specify %s envionment variable\n", envTestPort)
			return 1
		}

		ctx := context.Background()

		if err := os.Chdir("../"); err != nil {
			fmt.Fprintf(os.Stderr, "failed to change directory to the project root: %+v\n", err)
			return 1
		}

		cmd := exec.CommandContext(ctx, "make", "test-server")

		cmd.Env = append(
			os.Environ(),
			fmt.Sprintf("%s=%s", envPort, port),
			fmt.Sprintf("%s=%s", envTestPort, testPort),
		)

		if testing.Verbose() {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start the test server: %+v\n", err)
			return 1
		}

		conn, err := grpc.DialContext(
			ctx,
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

			termConn, err := net.Dial("tcp", fmt.Sprintf(":%s", testPort))
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to connect to the test server for termination: %+v\n", err)
				status = 1
			} else {
				err = termConn.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to close the connection to the test server for termination: %+v\n", err)
					status = 1
				}

				if err := cmd.Wait(); err != nil {
					fmt.Fprintf(os.Stderr, "failed to wait the test server: %+v\n", err)
					status = 1
				}
			}
		}()
		return m.Run()
	}())
}
