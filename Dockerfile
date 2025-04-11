FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ping-monitor

# Create the final image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

# Create app directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ping-monitor /app/
COPY --from=builder /app/embed/config.env /app/.env.example

# Create volume mount point for config
VOLUME ["/app/config"]

# Set environment variables
ENV TZ=Africa/Nairobi

# Run the application
CMD ["./ping-monitor"]
