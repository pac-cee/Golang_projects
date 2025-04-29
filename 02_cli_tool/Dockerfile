# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o cli_tool main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/cli_tool ./
CMD ["./cli_tool"]
