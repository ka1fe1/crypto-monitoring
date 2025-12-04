# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# -s: disable symbol table
# -w: disable DWARF generation
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o crypto-monitoring cmd/server/main.go

# Final stage
FROM scratch

WORKDIR /app

# Copy ca-certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /app/crypto-monitoring .

# Set config path environment variable
ENV CONFIG_PATH=/app/config/config.yaml

# Expose config volume
VOLUME ["/app/config"]

# Expose port
EXPOSE 8080

# Run the application
CMD ["./crypto-monitoring"]
