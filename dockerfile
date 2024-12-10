# Stage 1: Build the Go application
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go app (assuming the entry point is in the api directory)
RUN go build -o myapp ./api/main.go

# Stage 2: Create a minimal image
FROM alpine:latest

# Install necessary CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/myapp .

# Command to run the executable
CMD ["./myapp"]

