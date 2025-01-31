# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install dependencies (GCC & SQLite libraries required for CGO)
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Enable CGO for SQLite support
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the Go application
RUN go build -o main ./cmd/api/

# Stage 2: Create a lightweight container to run the app
FROM alpine:latest

# Install SQLite runtime dependencies
RUN apk add --no-cache sqlite

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Command to run the binary
CMD ["./main"]

