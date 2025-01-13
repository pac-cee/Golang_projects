.PHONY: build run test lint clean migrate-up migrate-down docker-up docker-down

# Build the application
build:
	go build -o bin/app

# Run the application
run:
	go run main.go

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/

# Database migrations
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/backend_master?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/backend_master?sslmode=disable" down

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Generate Swagger documentation
swagger:
	swag init -g internal/handler/handler.go

# Install dependencies
deps:
	go mod download

# Run development server with hot reload
dev:
	air

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
