.PHONY: up down test build run

up:
	docker-compose up --build -d

down:
	docker-compose down

test:
	go test -v ./...

# For running locally without docker (assuming DB is exposed on 5432)
run:
	DATABASE_URL="postgres://admin:password@localhost:5432/customer_registry?sslmode=disable" go run ./cmd/api/