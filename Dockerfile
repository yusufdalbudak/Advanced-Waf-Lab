# Multi-stage Dockerfile for production WAF
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wafd ./cmd/wafd

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 waf && \
    adduser -D -u 1000 -G waf waf

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/wafd .
COPY --from=builder /build/configs ./configs

# Create log directory
RUN mkdir -p /app/logs && \
    chown -R waf:waf /app

# Switch to non-root user
USER waf

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./wafd", "-config", "configs/waf.yaml"]

