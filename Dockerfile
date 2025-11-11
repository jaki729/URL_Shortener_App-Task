# Build stage
FROM golang:1.24 as builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/server

# Final runtime stage
FROM debian:stable-slim

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the service port
EXPOSE 8080

# Run the application
CMD ["./main"]