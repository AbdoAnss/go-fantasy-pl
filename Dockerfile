# Use the official Golang image to build the app
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the redis example
RUN go build -o redis-example ./examples/redis/main.go

# Use a minimal alpine image for the final container
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/redis-example .

# Run the binary
CMD ["./redis-example"]
