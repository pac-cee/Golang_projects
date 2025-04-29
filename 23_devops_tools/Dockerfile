# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o devops_tools main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/devops_tools ./
CMD ["./devops_tools"]
