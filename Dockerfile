# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26.2-alpine AS builder

WORKDIR /src

ARG TARGETOS
ARG TARGETARCH
ARG APP_NAME=mysql-mcp-server
ARG APP_DESC="A Go-based MCP server for MySQL."
ARG VERSION=develop
ARG BUILD=current
ARG COMMIT=none

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-$(go env GOARCH)} go build \
    -trimpath \
    -ldflags="-s -w \
    -X 'github.com/betterde/mysql-mcp-server/internal/build.Version=${VERSION}' \
    -X 'github.com/betterde/mysql-mcp-server/internal/build.Build=${BUILD}' \
    -X 'github.com/betterde/mysql-mcp-server/internal/build.Commit=${COMMIT}'" \
    -o /out/mysql-mcp-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S mysql-mcp-server \
    && adduser -S -G mysql-mcp-server mysql-mcp-server

WORKDIR /etc/mcp-servers/mysql

COPY --from=builder /out/mysql-mcp-server /usr/local/bin/mysql-mcp-server

USER mysql-mcp-server:mysql-mcp-server

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/mysql-mcp-server"]
CMD ["serve"]
