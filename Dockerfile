# Start with the Go base image
FROM golang:1.23.4

# Set environment variables
ENV GO111MODULE=on
ENV PORT=8080

# Create and set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Set the working directory to the location of main.go
WORKDIR /app/backend/cmd/server

# Build the Go application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
