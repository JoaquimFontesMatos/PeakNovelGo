# Stage 1: Build the Go application
FROM golang:1.23.4 AS builder

# Set environment variables
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Create and set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY backend/go.mod backend/go.sum ./

# Download Go dependencies
RUN go mod download

# Install Swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the application code
COPY . .

# Set the working directory to the location of main.go
WORKDIR /app/backend/cmd/server

# Generate Swagger documentation
RUN swag init -d ../../cmd/server,../../internal/controllers,../../internal/dtos,../../internal/models --parseDependency -o ../../docs

# Build the Go application with static linking
RUN go build -o main .

# Stage 2: Runtime environment
FROM debian:bullseye-slim

# Install Python and pip
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Copy the Go binary from the builder stage
COPY --from=builder /app/backend/cmd/server/main /app/main

# Copy the generated Swagger documentation
COPY --from=builder /app/docs /app/docs

# Copy Python requirements and the entire Python module
COPY backend/novel_updates_scraper /app/novel_updates_scraper
COPY backend/novel_updates_scraper/requirements.txt /app/requirements.txt

# Install Python dependencies
RUN pip3 install -r /app/requirements.txt

# Set PYTHONPATH so Python can locate the module
ENV PYTHONPATH=/app

# Set the working directory
WORKDIR /app

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
