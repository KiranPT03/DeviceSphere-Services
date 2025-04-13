# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy only the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum* ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations for size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Final stage
FROM alpine:latest
WORKDIR /app

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Copy configuration files
COPY --from=builder /app/pkg/config/application.yaml ./pkg/config/
COPY --from=builder /app/pkg/files/ ./pkg/files/

# Set the user to non-root
RUN adduser -D appuser
USER appuser

EXPOSE 3000
CMD ["./main"]
