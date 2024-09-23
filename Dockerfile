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

# Install ca-certificates for HTTPS requests (Docker API uses HTTPS)
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/bin/docodash .

# Copy the web folder with templates and static files
COPY --from=builder /app/web ./web

# Expose the port that the application will listen on
EXPOSE 8080

# Command to run the Go application
CMD ["./docodash"]