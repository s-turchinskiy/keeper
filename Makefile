.PHONY: test build clean lint deps

BIN_DIR = bin

define get_ldflags
    -X github.com/s-turchinskiy/keeper/internal/buildinfo.Version=$(get_version) \
    -X github.com/s-turchinskiy/keeper/internal/buildinfo.BuildTime=$(get_build_time) \
    -X github.com/s-turchinskiy/keeper/internal/buildinfo.Commit=$(get_commit)
endef

genproto:
	protoc --go_out=. --go-grpc_out=. internal/proto/api.proto --go_opt=default_api_level=API_OPAQUE

test:
	@go test -v ./...

test-integration:
	@go test -tags=integration ./internal/server/repository/ -v

fmt:
	@go fmt ./...

deps:
	go mod download
	go mod verify

build-client:
	@mkdir -p bin/
	@go build -ldflags="$(get_ldflags)" -o bin/keeperclient ./cmd/client

build-server:
	@mkdir -p bin
	@go build -ldflags="$(get_ldflags)" -o bin/keepersrv ./cmd/server

build-all: build-client build-server

lint:
	@golangci-lint run

clean:
	@rm -rf bin/

all: deps test build-all
