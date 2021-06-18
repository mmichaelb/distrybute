PROJECT_NAME=distrybute

GIT_VERSION=$(shell git describe --always)
GIT_BRANCH=$(shell git branch --show-current)

LD_FLAGS = -X main.GitVersion=${GIT_VERSION} -X main.GitBranch=${GIT_BRANCH}

OUTPUT_SUFFIX=$(go env GOEXE)

OUTPUT_PREFIX=./bin/${PROJECT_NAME}-${GIT_VERSION}

OUTPUT_FILE_ENDING=$(shell go env GOEXE)

# test go program
test:
	@go test -v ./...

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	@go build -ldflags '${LD_FLAGS}' -o "${OUTPUT_PREFIX}-${GOOS}-${GOARCH}${OUTPUT_FILE_ENDING}" ./cmd/distrybute/main.go

# installs and formats the project with the built-in Golang tool
install:
	@go install -ldflags '${LD_FLAGS}' ./cmd/distrybute/main.go