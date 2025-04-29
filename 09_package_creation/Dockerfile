# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o package_creation main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/package_creation ./
CMD ["./package_creation"]
