# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o tcp_server_client main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tcp_server_client ./
CMD ["./tcp_server_client"]
