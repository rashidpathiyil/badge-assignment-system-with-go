FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Use a small alpine image for the final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/server .

# Set environment variables
ENV GIN_MODE=release

# Expose the server port
EXPOSE 8080

# Run the server
CMD ["./server"] 
