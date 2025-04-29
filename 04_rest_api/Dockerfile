# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rest_api main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rest_api ./
CMD ["./rest_api"]
