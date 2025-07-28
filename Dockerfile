# Build stage
FROM golang:1.22-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ted ./cmd/ted

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/ted .

# Copy configuration and examples
COPY --from=builder /app/ted.config.yaml .
COPY --from=builder /app/examples ./examples

# Create output directory
RUN mkdir -p output

# Expose port for web dashboard
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./ted"] 