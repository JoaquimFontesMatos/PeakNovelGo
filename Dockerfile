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

# Copy the rest of the application code
COPY . .

# Set the working directory to the location of main.go
WORKDIR /app/backend/cmd/server

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

# Copy Python requirements and install dependencies
COPY backend/novel_updates_scraper/requirements.txt /app/requirements.txt
RUN pip3 install -r /app/requirements.txt

# Set the working directory
WORKDIR /app

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]