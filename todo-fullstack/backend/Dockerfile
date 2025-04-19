FROM golang:1.20
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
CMD ["go", "run", "main.go"]
