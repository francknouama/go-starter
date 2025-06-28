# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o verify-cli .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh verify-cli

# Set working directory
WORKDIR /home/verify-cli

# Copy binary from builder stage
COPY --from=builder /app/verify-cli .
COPY --from=builder /app/configs/ ./configs/

# Set ownership
RUN chown -R verify-cli:verify-cli /home/verify-cli

# Switch to non-root user
USER verify-cli

# Expose port (if needed)
# EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./verify-cli"]