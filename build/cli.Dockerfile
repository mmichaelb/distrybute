# build golang image
FROM golang:1.18.1-alpine3.15@sha256:42d35674864fbb577594b60b84ddfba1be52b4d4298c961b46ba95e9fb4712e8 AS builder
ARG build_git_branch=<unknown>
ARG build_git_tag=<unknown>
ARG build_git_commit_sha=<unknown>
LABEL stage=distrybute-builder
RUN apk update && apk add --no-cache git tzdata gcc libc-dev
WORKDIR $GOPATH/src/distrybute
COPY . .
RUN go mod verify
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags="-X github.com/mmichaelb/distrybute/internal/util.GitBranch=${build_git_branch} -X github.com/mmichaelb/distrybute/internal/util.GitTag=${build_git_tag} -X github.com/mmichaelb/distrybute/internal/util.GitCommitSha=${build_git_commit_sha} -w -s" \
     -o /go/bin/distrybute-cli ./cmd/distrybute-cli/main.go

# build real image
FROM scratch
ARG build_git_branch=<unknown>
ARG build_git_tag=<unknown>
ARG build_git_commit_sha=<unknown>
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/distrybute-cli /usr/bin/distrybute-cli
ENTRYPOINT ["/usr/bin/distrybute-cli"]

LABEL build_git_branch=${build_git_branch}
LABEL build_git_tag=${build_git_tag}
LABEL build_git_commit_sha=${build_git_commit_sha}
LABEL repository=https://github.com/mmichaelb/distrybute
