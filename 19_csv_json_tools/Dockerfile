# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o csv_json_tools main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/csv_json_tools ./
CMD ["./csv_json_tools"]
