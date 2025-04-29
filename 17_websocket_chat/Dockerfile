# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o websocket_chat main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/websocket_chat ./
CMD ["./websocket_chat"]
