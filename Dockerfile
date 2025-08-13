# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.24.6-alpine AS builder

WORKDIR /build

# Install build dependencies
# RUN apk add --no-cache git

ARG TARGETOS TARGETARCH

COPY go.mod go.sum ./

# Cache Go modules
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build for target OS/ARCH while running on the build platform to avoid emulation
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w" -o go-thermal-printer ./cmd/go-thermal-printer

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and curl
RUN apk --no-cache add ca-certificates curl

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /build/go-thermal-printer .

# Copy templates directory
COPY --from=builder /build/templates ./templates

# Copy config file
COPY --from=builder /build/config.example.toml ./config.toml

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Add appuser to dialout group
RUN addgroup appuser dialout

# Change ownership of app directory
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

ENV GIN_MODE=release

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl --fail --silent http://localhost:8080/health || exit 1

# Run the application
CMD ["./go-thermal-printer"]
