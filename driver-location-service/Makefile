.PHONY: up down test test-integration

up:
	docker compose up --build -d

down:
	docker compose down -v

test:
	go test ./...

test-integration:
	MONGO_URI=mongodb://localhost:27017 go test -tags=integration ./test/integration/...
