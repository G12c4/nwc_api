# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o nwc_app .

# Final stage
FROM alpine:latest

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/nwc_app .

# Copy .env file
COPY .env .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./nwc_app"]
