# build golang image
FROM golang:alpine AS builder
ARG build_git_branch=<unknown>
ARG build_git_tag=<unknown>
ARG build_git_commit_sha=<unknown>
LABEL stage=distrybute-builder
RUN apk update && apk add --no-cache git tzdata gcc libc-dev
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/distrybute
COPY . .
RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -trimpath \
    -ldflags="-X github.com/mmichaelb/distrybute/internal/app.GitBranch=${build_git_branch} -X github.com/mmichaelb/distrybute/internal/app.GitTag=${build_git_tag} -X github.com/mmichaelb/distrybute/internal/app.GitCommitSha=${build_git_commit_sha} -w -s" \
     -o /go/bin/distrybute ./cmd/distrybute/main.go

# build real image
FROM scratch
ARG build_git_branch=<unknown>
ARG build_git_tag=<unknown>
ARG build_git_commit_sha=<unknown>
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/distrybute /usr/bin/distrybute
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER appuser:appuser
#RUN chmod +x /app/distrybute/bin/distrybute
ENTRYPOINT ["/usr/bin/distrybute"]

LABEL build_git_branch=${build_git_branch}
LABEL build_git_tag=${build_git_tag}
LABEL build_git_commit_sha=${build_git_commit_sha}
LABEL repository=https://github.com/mmichaelb/distrybute
