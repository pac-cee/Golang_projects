# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o context_middleware main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/context_middleware ./
CMD ["./context_middleware"]
