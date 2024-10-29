# Step 1: Download JS dependencies
FROM node:23.1.0-alpine3.20 AS js-dependencies

# Necessary to avoid https://stackoverflow.com/a/65443098/106918
WORKDIR /js

# Install the required dependencies
RUN npm install htmx.org@1.9.12
RUN npm install mermaid@11.3.0

# Step 2: Build the Go binary
FROM golang:1.23.2-alpine3.20 AS builder

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
RUN go build -o bin/laebel ./cmd/laebel

# Step 3: Create a minimal image to run the application
FROM alpine:3.20.3 AS laebel

# Expose the port that the application will listen on
ENV PORT=8000
EXPOSE $PORT

# Configure a healthcheck
HEALTHCHECK \
    --timeout=5s --start-period=3s --start-interval=3s \
    CMD wget -q -S -O - http://localhost:$PORT/ || exit 1

# Command to run the Go application
CMD ["./laebel"]

# Set documentation labels
LABEL org.opencontainers.image.title="Laebel" \
    org.opencontainers.image.description="Automatic README-style documentation site for your Docker Compose project." \
    org.opencontainers.image.authors="Henrik Jernevad <henrik@jernevad.se>" \
    net.henko.laebel.group="Documentation" \
    net.henko.laebel.hidden="true" \
    net.henko.laebel.port.$PORT.description="The port where this documentation site is served."

# Set the working directory inside the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/bin/laebel .

# Copy the web folder with templates and static files
COPY --from=builder /app/web ./web

# Copy Mermaid dependencies
COPY --from=js-dependencies /js/node_modules/mermaid/dist/mermaid.esm.min.mjs ./web/static/js/mermaid.esm.min.mjs
COPY --from=js-dependencies /js/node_modules/mermaid/dist/chunks/mermaid.esm.min/* ./web/static/js/chunks/mermaid.esm.min/

# Copy HTMX dependencies
COPY --from=js-dependencies /js/node_modules/htmx.org/dist/htmx.min.js ./web/static/js/htmx.min.js
COPY --from=js-dependencies /js/node_modules/htmx.org/dist/ext/sse.js ./web/static/js/sse.js