FROM golang:1.25-alpine

# Install additional tools you might need
RUN apk add --no-cache \
    git \
    curl \
    nano \
    make

# Set working directory
WORKDIR /workspace

# Pre-download dependencies to speed up development
RUN go mod download || true
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest