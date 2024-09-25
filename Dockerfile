# Step 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install Git and required dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project
COPY cmd/ ./cmd/
COPY web/ ./web/

# Build the Go binary
RUN go build -o bin/docodash ./cmd/docodash

# Step 2: Create a minimal image to run the application
FROM alpine:latest AS docodash

# Expose the port that the application will listen on
ENV PORT=8080
EXPOSE $PORT

# Configure a healthcheck
HEALTHCHECK \
    --timeout=5s --start-period=3s --start-interval=3s \
    CMD wget -q -S -O - http://localhost:$PORT/ || exit 1

# Command to run the Go application
CMD ["./docodash"]

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/bin/docodash .

# Copy the web folder with templates and static files
COPY --from=builder /app/web ./web
