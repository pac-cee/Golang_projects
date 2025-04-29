# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o image_processing main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/image_processing ./
CMD ["./image_processing"]
