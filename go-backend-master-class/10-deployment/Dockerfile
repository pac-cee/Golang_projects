# Start from the official Go image
FROM golang:1.21-alpine

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=0 /app/main .

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
