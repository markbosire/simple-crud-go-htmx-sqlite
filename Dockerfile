# Start from the official Golang image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o taskmanager .

# Start from a minimal Alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/taskmanager .

# Copy static files
COPY ./static ./static

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./taskmanager"]