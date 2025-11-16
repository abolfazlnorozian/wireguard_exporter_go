# Multi-stage build for smaller final image
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' \
    -o wireguard_exporter ./cmd/wireguard_exporter

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    wireguard-tools \
    iptables \
    iproute2

# Create non-root user (though we'll need root for WireGuard access)
RUN addgroup -g 1000 exporter && \
    adduser -D -u 1000 -G exporter exporter

# Copy binary from builder
COPY --from=builder /build/wireguard_exporter /usr/local/bin/wireguard_exporter

# Make it executable
RUN chmod +x /usr/local/bin/wireguard_exporter

# Expose metrics port
EXPOSE 9586

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:9586/metrics || exit 1

# Run as root (required for WireGuard access)
USER root

ENTRYPOINT ["/usr/local/bin/wireguard_exporter"]
CMD ["-listen-address", ":9586"]

