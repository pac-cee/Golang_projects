# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o web_scraper main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/web_scraper ./
CMD ["./web_scraper"]
