# Start from the official Go image
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o task-manager

# Final stage
FROM alpine:latest

# Install sqlite runtime dependencies
RUN apk add --no-cache sqlite

# Set working directory
WORKDIR /app

# Create data directory
RUN mkdir -p /app/data

# Copy the built binary
COPY --from=builder /app/task-manager .

# Copy templates and static files
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose port
EXPOSE 8080

# Set default database path inside the container
ENV DATABASE_PATH=/app/data/tasks.db

# Command to run the executable
CMD ["./task-manager"]