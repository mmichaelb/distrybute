PROJECT_NAME=distrybute

GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG=$(shell git describe --tags --always)
GIT_COMMIT_SHA=$(shell git rev-parse HEAD)

LD_FLAGS = -X github.com/mmichaelb/distrybute/internal/util.GitBranch=${GIT_BRANCH} -X github.com/mmichaelb/distrybute/internal/util.GitTag=${GIT_TAG} -X github.com/mmichaelb/distrybute/internal/util.GitCommitSha=${GIT_COMMIT_SHA}

OUTPUT_SUFFIX=$(shell go env GOEXE)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

OUTPUT_PREFIX=./bin/${PROJECT_NAME}-${GIT_TAG}

OUTPUT_FILE_ENDING=$(shell go env GOEXE)

# unit test go program
.PHONY: unit-test
unit-test:
	@go test -race -coverpkg=./pkg/... -coverprofile=coverage_unit.txt -covermode=atomic -v ./...
	@sed -i '/github.com\/mmichaelb\/distrybute\/pkg\/mocks\//d' coverage_unit.txt

.PHONY: postgres-minio-integration-test
postgres-minio-integration-test:
	@POSTGRES_MINIO_INTEGRATION_TEST=true go test -race -coverprofile=coverage_integration.txt -covermode=atomic -v -run Test_PostgresMinio_Service ./...
	@sed -i '/github.com\/mmichaelb\/distrybute\/pkg\/mocks\//d' coverage_integration.txt

# builds and formats the project with the built-in Golang tool
.PHONY: build
build:
	@go build -trimpath -ldflags '${LD_FLAGS}' -o "${OUTPUT_PREFIX}-${GOOS}-${GOARCH}${OUTPUT_FILE_ENDING}" ./cmd/distrybute/main.go

.PHONY: build-cli
build-cli:
	@go build -trimpath -ldflags '${LD_FLAGS}' -o "./bin/${PROJECT_NAME}-cli-${GIT_TAG}-${GOOS}-${GOARCH}${OUTPUT_FILE_ENDING}" ./cmd/distrybute-cli/main.go

# installs and formats the project with the built-in Golang tool
.PHONY: install
install:
	@go install -trimpath -ldflags '${LD_FLAGS}' ./cmd/distrybute/main.go

.PHONY: swagger
swagger:
	@swag init --parseInternal --generatedTime -g pkg/rest/controller/router.go

.PHONY: swagger-format
swagger-format:
	@swag fmt -g router.go -d pkg/rest/controller/

.PHONY: deps
deps:
	@go mod download

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: mockery
mockery:
	@mockery --dir pkg/ --name ".*" --keeptree

.PHONY: docker-build
docker-build: vendor
	@docker build -f "build/Dockerfile" \
		-t ghcr.io/mmichaelb/distrybute:${GIT_TAG} -t ghcr.io/mmichaelb/distrybute:latest \
		-t mmichaelb/distrybute:${GIT_TAG} -t mmichaelb/distrybute:latest \
		--build-arg build_git_branch=${GIT_BRANCH} --build-arg build_git_tag=${GIT_TAG} --build-arg build_git_commit_sha=${GIT_COMMIT_SHA} \
		.

.PHONY: docker-cross-platform-buildx-push
docker-cross-platform-buildx-push: vendor
	@docker buildx build --platform linux/amd64,linux/arm/v7,linux/arm64/v8 -f "build/Dockerfile" --push \
		-t ghcr.io/mmichaelb/distrybute:${GIT_TAG} -t ghcr.io/mmichaelb/distrybute:latest \
		-t mmichaelb/distrybute:${GIT_TAG} -t mmichaelb/distrybute:latest \
		--build-arg build_git_branch=${GIT_BRANCH} --build-arg build_git_tag=${GIT_TAG} --build-arg build_git_commit_sha=${GIT_COMMIT_SHA} \
		.

.PHONY: docker-build-cli
docker-build-cli: vendor
	@docker build -f "build/Dockerfile-cli" \
		-t ghcr.io/mmichaelb/distrybute-cli:${GIT_TAG} -t ghcr.io/mmichaelb/distrybute-cli:latest \
		-t mmichaelb/distrybute-cli:${GIT_TAG} -t mmichaelb/distrybute-cli:latest \
		--build-arg build_git_branch=${GIT_BRANCH} --build-arg build_git_tag=${GIT_TAG} --build-arg build_git_commit_sha=${GIT_COMMIT_SHA} \
		.

.PHONY: docker-cross-platform-buildx-push-cli
docker-cross-platform-buildx-push-cli: vendor
	@docker buildx build --platform linux/amd64,linux/arm/v7,linux/arm64/v8 -f "build/Dockerfile-cli" --push \
		-t ghcr.io/mmichaelb/distrybute-cli:${GIT_TAG} -t ghcr.io/mmichaelb/distrybute-cli:latest \
		-t mmichaelb/distrybute-cli:${GIT_TAG} -t mmichaelb/distrybute-cli:latest \
		--build-arg build_git_branch=${GIT_BRANCH} --build-arg build_git_tag=${GIT_TAG} --build-arg build_git_commit_sha=${GIT_COMMIT_SHA} \
		.

.PHONY: remove-docker-helper-images
remove-docker-helper-images:
	@docker image prune --force --filter label=stage=distrybute-builder
