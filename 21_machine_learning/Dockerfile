# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o machine_learning main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/machine_learning ./
CMD ["./machine_learning"]
