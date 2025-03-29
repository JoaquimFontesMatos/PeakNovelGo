# Stage 1: Build the Go application and generate docs
FROM golang:1.23.4 AS builder

# Set environment variables
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Set working directory
WORKDIR /app

# Copy only the module files to leverage caching for dependency download
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Install Swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the source code
COPY . .

# Change working directory to backend and ensure docs directory exists
WORKDIR /app/backend
RUN mkdir -p docs

# Generate Swagger documentation (using corrected paths)
RUN swag init -d ./cmd/server,./internal/controllers,./internal/dtos,./internal/models --parseDependency -o ./docs && ls -lah ./docs

# Change directory to where the main Go binary is built and build it with static linking
WORKDIR /app/backend/cmd/server
RUN go build -o main .

# Stage 2: Runtime environment
FROM python:3.13.1-slim

# Copy Go binary and generated Swagger docs from the builder stage
COPY --from=builder /app/backend/cmd/server/main /app/main
COPY --from=builder /app/backend/docs /app/docs

# Optimize Python dependency installation:
# 1. Copy only the requirements first and install (to cache pip install layer)
COPY backend/novel_updates_scraper/requirements.txt /app/requirements.txt
RUN pip3 install --upgrade pip && \
    pip3 install --no-cache-dir -r /app/requirements.txt

# 2. Then copy the rest of the Python module code
COPY backend/novel_updates_scraper /app/novel_updates_scraper

# Combine Playwright installation commands into one layer
RUN playwright install && playwright install-deps

# Set PYTHONPATH so Python can locate your module
ENV PYTHONPATH=/app

# Set working directory and expose port
WORKDIR /app
EXPOSE 8080

# Run the application (Go binary)
CMD ["./main"]
