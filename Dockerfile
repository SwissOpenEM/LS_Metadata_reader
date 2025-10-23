# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build static binary
ARG VERSION=0.0.0
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -C cmd/LS_Metadata_reader \
    -ldflags="-s -w -X 'main.version=${VERSION}'" \
    -o /src/LS_Metadata_reader \
    .

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /src/LS_Metadata_reader /usr/local/bin/LS_Metadata_reader

# Run as non-root
RUN adduser -D -g '' appuser && chown appuser:appuser /usr/local/bin/LS_Metadata_reader
USER appuser

ENTRYPOINT ["/usr/local/bin/LS_Metadata_reader"]
