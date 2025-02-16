# Use the official Alpine image as a base
FROM golang:1.24-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the Go mod file
COPY go.mod ./

# Download the dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port
EXPOSE 3000

# Run the command to start the application
CMD ["./main"]
