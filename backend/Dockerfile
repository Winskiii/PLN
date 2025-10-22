# syntax=docker/dockerfile:1.7-labs

FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache \
    apk add --no-cache git && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
ENV GIN_MODE=release
USER nonroot:nonroot
WORKDIR /
COPY --chown=nonroot:nonroot --from=builder /out/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
