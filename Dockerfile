# Stage 1: Build the application
FROM oven/bun:1.2-slim AS web-builder

# Set the working directory
WORKDIR /web

# Copy the React application source code
COPY cmd/web/ .

# Install dependencies and build the React application
RUN bun install && bun run build

# Stage 2: Build the Go application
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set the working directory
WORKDIR /app

# Copy the Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Enable CGO and build the Go application
ENV CGO_ENABLED=1
RUN go build -o api ./cmd/api/

# Stage 3: Create a minimal runtime image
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk add --no-cache sqlite-libs

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/api .

# Copy the built React static files from the web-builder stage
COPY --from=web-builder /web/dist ./static
COPY --from=web-builder /web/src/logo.svg ./static

# Expose the data volume
VOLUME /app/data

# Expose the port for the API
EXPOSE 3000

# Set the NODE_ENV environment variable to production
ENV NODE_ENV=production

# Command to run the API
CMD ["./api"]