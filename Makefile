PORT := 9090
TEST_PORT := 9091
GO_PROTOC_VERSION := 1.3.1

.PHONY: test
test:
	@PORT=${PORT} TEST_PORT=${TEST_PORT} go test -v -race -count=1 ./e2e

.PHONY: test-server
test-server:
	@go test -race -count=1 . -covermode=atomic -coverprofile=coverage.out -run='^TestRunMain$$'

.PHONY: pb
pb:
	@docker run --rm --volume $(shell pwd):/go/src/github.com/110y/go-e2e-example 110y/go-protoc:${GO_PROTOC_VERSION} protoc \
		-I /go/src/github.com/110y/go-e2e-example/server/pb/ \
		--go_out=plugins=grpc:/go/src/github.com/110y/go-e2e-example/server/pb \
		/go/src/github.com/110y/go-e2e-example/server/pb/server.proto
