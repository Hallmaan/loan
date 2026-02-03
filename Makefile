.PHONY: build run test clean migrate migrate-down docker-up docker-down

# Application
APP_NAME=amartha
MAIN_PATH=./cmd/api

# Database
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/amartha?sslmode=disable

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -rf uploads/

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database migrations
migrate:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Development
dev: docker-up
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 3
	$(MAKE) migrate
	$(MAKE) run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install migrate tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
