# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-boilerplate .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/gin-boilerplate .

# Copy configuration files
COPY --from=builder /app/config.prod.yaml ./config.prod.yaml

# Expose port
EXPOSE 8080

# Run the application
CMD ["./gin-boilerplate", "-env=prod"]
