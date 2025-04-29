# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o advanced_patterns main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/advanced_patterns ./
CMD ["./advanced_patterns"]
