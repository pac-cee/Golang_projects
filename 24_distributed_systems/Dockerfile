# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o distributed_systems main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/distributed_systems ./
CMD ["./distributed_systems"]
