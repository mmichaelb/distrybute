PROJECT_NAME=distrybute

GIT_VERSION=$(shell git describe --always)
GIT_BRANCH=$(shell git branch --show-current)

LD_FLAGS = -X github.com/mmichaelb/distrybute/internal/app.GitVersion=${GIT_VERSION} -X github.com/mmichaelb/distrybute/internal/app.GitBranch=${GIT_BRANCH}

OUTPUT_SUFFIX=$(shell go env GOEXE)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

OUTPUT_PREFIX=./bin/${PROJECT_NAME}-${GIT_VERSION}

OUTPUT_FILE_ENDING=$(shell go env GOEXE)

# unit test go program
unit-test:
	@go test -race -coverpkg=./pkg/... -coverprofile=coverage.txt -covermode=atomic -v ./...
	@sed -i '/github.com\/mmichaelb\/distrybute\/pkg\/mocks\//d' coverage.txt

postgres-minio-integration-test:
	@POSTGRES_MINIO_INTEGRATION_TEST=true go test -race -coverprofile=coverage.txt -covermode=atomic -v -run Test_PostgresMinio_Service ./...

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	@go build -trimpath -ldflags '${LD_FLAGS}' -o "${OUTPUT_PREFIX}-${GOOS}-${GOARCH}${OUTPUT_FILE_ENDING}" ./cmd/distrybute/main.go

.PHONY: build-cli
build-cli:
	@go build -trimpath -ldflags '${LD_FLAGS}' -o "./bin/${PROJECT_NAME}-cli-${GIT_VERSION}-${GOOS}-${GOARCH}${OUTPUT_FILE_ENDING}" ./cmd/distrybute-cli/main.go

# installs and formats the project with the built-in Golang tool
install:
	@go install -trimpath -ldflags '${LD_FLAGS}' ./cmd/distrybute/main.go

swagger:
	@swag init --parseInternal --generatedTime -g pkg/rest/controller/router.go

swagger-format:
	@swag fmt -g router.go -d pkg/rest/controller/

deps:
	@go mod download

mockery:
	@mockery --dir pkg/ --name ".*" --keeptree
