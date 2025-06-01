# # Use the official Golang image as a base
# FROM golang:1.24-alpine

# # Install build dependencies for CGO and SQLite
# RUN apk add --no-cache gcc musl-dev sqlite-dev

# # Set the working directory
# WORKDIR /app

# # Copy the Go modules and source code
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . .

# # Build the application
# RUN go build -o api ./cmd/api/main.go

# # Expose the data volume
# VOLUME /app/data

# # Expose the port for the API
# EXPOSE 3000

# # Command to run the API
# CMD ["./api"]

# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory
WORKDIR /app

# Copy the Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Enable CGO and build the application
ENV CGO_ENABLED=1
RUN go build -o api ./cmd/api/main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk add --no-cache sqlite-libs

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/api .

# Expose the data volume
VOLUME /app/data

# Expose the port for the API
EXPOSE 3000

# Command to run the API
CMD ["./api"]