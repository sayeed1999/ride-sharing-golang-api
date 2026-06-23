### ─────────────────────────────────────────────
### Multi-stage Dockerfile for api
### Builds a minimal and production-ready Go container image
### ─────────────────────────────────────────────

## ────────────────
## Stage 1: Builder
## ────────────────
# Use the official lightweight Go image with Alpine as the base
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /src

# Install git (required for Go modules that fetch from private/public repos)
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download and cache Go dependencies
RUN go mod download

# Copy the entire project source into the container
COPY . .

# Build the Go binary for Linux, statically linked (no CGO dependencies)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -o /out/app ./cmd/app


## ────────────────
## Stage 2: Runtime
## ────────────────
# Use a fresh, minimal Alpine image for the final runtime
FROM alpine:3.22

# Add CA certificates to enable HTTPS requests from the app
RUN apk add --no-cache ca-certificates

# Copy only the compiled Go binary from the builder stage
COPY --from=builder /out/app /usr/local/bin/app

# Document the port the service listens on (informational)
EXPOSE 8080

# Set the container entrypoint to start the Go API
ENTRYPOINT ["/usr/local/bin/app"]
